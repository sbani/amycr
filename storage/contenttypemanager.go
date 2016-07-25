package storage

import "github.com/sbani/amycr/contenttype"

// ContentTypeManager interface to store contenttypes in storage
type ContentTypeManager interface {
	// Put inserts or creates a new contenttype
	Put(contenttype.ContentType) error

	// Get return a single contenttype entry
	Get(key string) (contenttype.ContentType, error)

	// Delete removes a single contenttype
	Delete(contenttype.ContentType) error

	// FindAll returns the complete list of all content types
	FindAll() ([]contenttype.ContentType, error)
}
