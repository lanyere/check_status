package check_status

import (
	"context"
	"log"
)

// MockLogger is a mock implementation of the Logger interface for testing purposes.
type MockLogger struct {
	SavedMessages []string
}

func (l *MockLogger) Save(message string, id string, options ...any) {
	log.Printf("Message: %s | ID: %s | Options: %v", message, id, options)
}

func (l *MockLogger) SaveMessage(message string) {
	log.Println(message)
}

// MockDatabase is a mock implementation of the Database interface for testing purposes.
type MockDatabase struct {
	Data map[string]*Transaction
}

func (db *MockDatabase) GetID(ctx context.Context, id string) (*Transaction, error) {
	return &Transaction{
		ID:     "123",
		Status: "incomplete",
	}, nil
}

func (db *MockDatabase) Update(ctx context.Context, tr *Transaction) error {
	return nil
}
