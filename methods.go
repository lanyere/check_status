package check_status

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (p *poll) verifyOptions() error {
	if len(p.config.Providers) <= 0 {
		return fmt.Errorf("at least one provider must be specified")
	}

	if p.config.Database == nil {
		return fmt.Errorf("db must be specified")
	}

	if p.config.Interval == 0 {
		p.config.Interval = defaultInterval
	}

	return nil
}

func (p *poll) startTicker() {
	ticker := time.NewTicker(p.config.Interval)
	defer ticker.Stop()

	p.config.Log.SaveMessage("Ticker is running")

	for {
		select {
		case <-ticker.C:
			if !p.status {
				p.config.Log.SaveMessage("Ticker is stopped")
				ticker.Stop()
				return
			}

			for _, provider := range p.config.Providers {
				start := time.Now()

				resp, err := http.Get(provider)
				if err != nil {
					p.config.Log.SaveMessage("can't reach the provider: " + provider)
					errorCount.WithLabelValues(provider).Inc()
					continue
				}

				requestDuration.WithLabelValues(provider).Observe(time.Since(start).Seconds())
				requestCount.WithLabelValues(provider).Inc()

				closeBody := func() {
					if err = resp.Body.Close(); err != nil {
						p.config.Log.SaveMessage("can't close response body: " + err.Error())
						errorCount.WithLabelValues(provider).Inc()
						return
					}
				}

				tr := &Callback{}
				if err = json.NewDecoder(resp.Body).Decode(&tr); err != nil {
					p.config.Log.SaveMessage("can't unmarshal the body from provider: " + provider)
					closeBody()
					continue
				}

				if len(tr.Transactions) <= 0 {
					p.config.Log.SaveMessage("no new transaction statuses from provider: " + provider)
					closeBody()
					continue
				}

				go p.handleTransactions(context.TODO(), tr.Transactions, provider)
				closeBody()
			}
		}
	}
}

func (p *poll) handleTransactions(ctx context.Context, trs []Transaction, provider string) {
	for i := range trs {
		tr := trs[i]

		dbTR, err := p.config.Database.GetID(ctx, tr.ID)
		if err != nil {
			p.config.Log.Save("transaction is not found due to "+err.Error(), tr.ID)
			continue
		}

		if dbTR == nil {
			p.config.Log.Save("transaction is not found", tr.ID)
			continue
		}

		if dbTR.Status == tr.Status {
			p.config.Log.Save("transaction has the same status", tr.ID)
			continue
		}

		dbTR.Status = trs[i].Status
		if err = p.config.Database.Update(ctx, dbTR); err != nil {
			p.config.Log.Save("transaction is not updated due to:"+err.Error(), tr.ID)
			return
		}

		p.data.Store(dbTR.ID, dbTR.Status)
		transactionsProcessed.WithLabelValues(provider).Inc()
	}
}
