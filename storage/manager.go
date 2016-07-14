package storage

// Manager interface for storage
type Manager interface {
	// PutRecord inserts or creates a new record
	PutRecord(r *Record) error

	// GetRecord return a single record entry
	GetRecord(contentType, key []byte) (Record, error)

	// DeleteRecord removes a single record
	DeleteRecord(r Record) error

	// FindNextRevision searches the an older revision than the current
	FindNextRevision(r *Record) (Record, error)

	// PutContentType updates or creates a content type by key
	PutContentType(c *ContentType) error

	// GetContentType a single entry
	GetContentType(key []byte) (ContentType, error)

	// DeleteContentType removes a content type from storage
	DeleteContentType(c ContentType) error
}
