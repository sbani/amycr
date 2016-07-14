package storage

import "time"

// Manager interface for storage
type Manager interface {
	// CreateRecord inserts a new record
	CreateRecord(contentType, content []byte) (Record, error)

	// CreateRecordWithkey inserts a new record but adds a fixed key (e.g.)
	CreateRecordWithkey(contentType, key, content []byte) (Record, error)

	// UpdateRecord updates record to a new version
	UpdateRecord(contentType, key, content []byte) (Record, error)

	// GetContant a single entry
	GetRecord(contentType, key []byte) (Record, error)

	// HasRecord checks if a record exists in db
	HasRecord(contentType, key []byte) bool

	// DeleteRecord removes a single record
	DeleteRecord(contentType, key []byte) error

	// Get the recods to a specific time
	GetRecordRevision(key []byte, t time.Time) (Record, error)

	// CreateContentType adds a new content type
	CreateContentType(key, content []byte) error

	// UpdateContentType updates a content type by key
	UpdateContentType(key, content []byte) error

	// GetContentType a single entry
	GetContentType(key []byte) ([]byte, error)

	// DeleteContentType removes a content type from storage
	DeleteContentType(key []byte) error
}
