package storage

import "time"

// Manager interface for storage
type Manager interface {
	// CreateRecord inserts a new record
	CreateRecord(contentType, content []byte) (Record, error)

	// CreateRecordWithID inserts a new record but adds a fixed key (e.g.)
	CreateRecordWithID(contentType, id, content []byte) (Record, error)

	// UpdateRecord updates record to a new version
	UpdateRecord(contentType, id, content []byte) (Record, error)

	// GetContant a single entry
	GetRecord(contentType, id []byte) (Record, error)

	// HasRecord checks if a record exists in db
	HasRecord(contentType, id []byte) bool

	// DeleteRecord removes a single record
	DeleteRecord(contentType, id []byte) error

	// Get the recods to a specific time
	GetRecordRevision(id []byte, t time.Time) (Record, error)

	// CreateContentType adds a new content type
	CreateContentType(id, content []byte) error

	// GetContentType a single entry
	GetContentType(id []byte) ([]byte, error)

	// DeleteContentType removes a content type from storage
	DeleteContentType(id []byte) error
}
