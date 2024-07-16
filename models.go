package check_status

import (
	"sync"
	"time"
)

const defaultInterval = time.Minute // ¯\_(ツ)_/¯

type Poller interface {
	StartPolling()
	StopPolling()
	GetTransactionStatus(id string) (Status, error)
}

type poll struct {
	config *Config
	status bool
	data   sync.Map
}

type Config struct {
	Providers []string
	Interval  time.Duration
	Log       Logger
	Database  Database
}

type Status string

const IncompleteStatus = Status("incomplete")
const CompleteStatus = Status("complete")

type Callback struct {
	Transactions []Transaction `json:"transactions,omitempty"`
}

type Transaction struct {
	ID     string `json:"id,omitempty"`
	Status Status `json:"status,omitempty" type:"Status"`
	// there is more fields, yet I'll use just these in the example
}
