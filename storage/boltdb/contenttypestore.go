package boltdb

import (
	"github.com/asdine/storm"
	"github.com/sbani/gcr/contenttype"
)

// ContentTypeStore represents the content type store
type ContentTypeStore struct {
	db *storm.DB
}

// NewContentTypeStore is the factory for ContentTypeStore
func NewContentTypeStore(db *storm.DB) *ContentTypeStore {
	return &ContentTypeStore{
		db: db,
	}
}

// Put inserts or creates a new contenttype
func (s *ContentTypeStore) Put(c *contenttype.ContentType) error {
	return s.db.Save(c)
}

// Get return a single contenttype entry
func (s *ContentTypeStore) Get(key string) (contenttype.ContentType, error) {
	var ct contenttype.ContentType
	err := s.db.One("Key", key, &ct)

	return ct, err
}

// Delete removes a single contenttype
func (s *ContentTypeStore) Delete(c *contenttype.ContentType) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	// Remove content type itself
	err = tx.Remove(c)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove record buckets
	_ = tx.Drop(c.Key)

	return tx.Commit()
}

// FindAll returns a list of all content types
func (s *ContentTypeStore) FindAll() ([]*contenttype.ContentType, error) {
	var ct []*contenttype.ContentType
	err := s.db.All(&ct)

	return ct, err
}
