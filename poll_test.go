package check_status

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"
)

var port = "8080"

func TestPolling(t *testing.T) {
	mockLogger := &MockLogger{}
	mockDB := &MockDatabase{
		Data: make(map[string]*Transaction),
	}

	p := NewPoller(&Config{
		Providers: []string{"http://localhost:" + port + "/tr"},
		Interval:  2,
		Log:       mockLogger,
		Database:  mockDB,
	})

	t.Run("TestPollingProcess", func(t *testing.T) {
		go mockServer(port)

		p.StartPolling()
		time.Sleep(6 * time.Millisecond)
		p.StopPolling()

		expectedTransactionID := "123"
		expectedStatus := CompleteStatus

		status, err := p.GetTransactionStatus(expectedTransactionID)
		if err != nil {
			log.Println(err)
			return
		}

		if status != expectedStatus {
			t.Errorf("Expected to retrieve transaction %s with Status %s, but got %s", expectedTransactionID, expectedStatus, status)
			return
		}

		t.Logf("Expected to retrieve transaction with Status %s got correct status", expectedStatus)
	})

	t.Run("TestMetricsExposure", func(t *testing.T) {
		go ExposeMetrics("/metrics", "2112")

		req, err := http.NewRequest("GET", "http://localhost:2112/metrics", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
			return
		}

		t.Logf("Expected status code 200, got %d", resp.StatusCode)
	})

	t.Run("TestMetrics", func(t *testing.T) {
		metrics, err := CollectMetrics()
		if err != nil {
			t.Fatalf("Failed to collect metrics: %v", err)
			return
		}

		if metrics == nil {
			t.Error("Expected metrics, got nil")
		}

		for _, metric := range metrics {
			m := metric.GetMetric()
			t.Logf("Metric: %v", m)
			return
		}

	})
}

func mockServer(port string) {
	http.HandleFunc("/tr", func(w http.ResponseWriter, r *http.Request) {
		transaction := &Callback{
			Transactions: []Transaction{
				{
					ID:     "123",
					Status: CompleteStatus,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
