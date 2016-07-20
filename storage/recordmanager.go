package storage

import "github.com/sbani/gcr/record"

// RecordManager interface to store records in storage
type RecordManager interface {
	// Put inserts or creates a new record
	Put(*record.Record) error

	// Get return a single record entry
	Get(contentType, key string) (record.Record, error)

	// Delete removes a single record
	Delete(*record.Record) error
}
