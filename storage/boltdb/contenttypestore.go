package boltdb

import (
	"github.com/asdine/storm"
	"github.com/hashicorp/golang-lru"
	"github.com/sbani/amycr/contenttype"
)

const (
	// CacheSize gives the number of elements in LRU Cache
	CacheSize = 100
)

// ContentTypeStore represents the content type store
type ContentTypeStore struct {
	db    *storm.DB
	cache *lru.Cache
}

// NewContentTypeStore is the factory for ContentTypeStore
func NewContentTypeStore(db *storm.DB) *ContentTypeStore {
	c, _ := lru.New(CacheSize)

	return &ContentTypeStore{
		db:    db,
		cache: c,
	}
}

// Put inserts or creates a new contenttype
func (s *ContentTypeStore) Put(c contenttype.ContentType) error {
	err := s.db.Save(&c)
	if err != nil {
		return err
	}

	s.cache.Add(c.Key, c)

	return nil
}

// Get return a single contenttype entry
func (s *ContentTypeStore) Get(key string) (contenttype.ContentType, error) {
	var c contenttype.ContentType
	if v, ok := s.cache.Get(key); ok {
		c = v.(contenttype.ContentType)
		return c, nil
	}

	err := s.db.One("Key", key, &c)

	s.cache.Add(c.Key, c)

	return c, err
}

// Delete removes a single contenttype
func (s *ContentTypeStore) Delete(c contenttype.ContentType) error {
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

	s.cache.Remove(c.Key)

	return tx.Commit()
}

// FindAll returns a list of all content types
func (s *ContentTypeStore) FindAll() ([]contenttype.ContentType, error) {
	var ct []contenttype.ContentType
	err := s.db.All(&ct)

	return ct, err
}
