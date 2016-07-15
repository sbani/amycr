package storage

// Manager interface for storage
type Manager interface {
	// PutRecord inserts or creates a new record
	PutRecord(*Record) error

	// GetRecord return a single record entry
	GetRecord([]byte, []byte) (Record, error)

	// DeleteRecord removes a single record
	DeleteRecord(Record) error

	// FindNextRevision searches the an older revision than the current
	FindNextRevision(*Record) (Record, error)

	// PutContentType updates or creates a content type by key
	PutContentType(*ContentType) error

	// GetContentType a single entry
	GetContentType([]byte) (ContentType, error)

	// ListContentTypes returns a slice of content types
	ListContentTypes() ([]ContentType, error)

	// DeleteContentType removes a content type from storage
	DeleteContentType(ContentType) error
}
