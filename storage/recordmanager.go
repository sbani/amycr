package storage

import "github.com/sbani/amycr/record"

// RecordManager interface to store records in storage
type RecordManager interface {
	// Put inserts or creates a new record
	Put(record.Record) error

	// Get return a single record entry
	Get(contentType, key string) (record.Record, error)

	// GetRevisions returns a list of revisions for a given record
	GetRevisions(record.Record) ([]record.Revision, error)

	// Delete removes a single record
	Delete(record.Record) error
}
