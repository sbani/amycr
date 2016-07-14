package storage

import (
	"bytes"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/evanphx/json-patch"
	"github.com/pkg/errors"
)

var (
	historyBucket     = []byte("history")
	contentTypeBucket = []byte("contenttype")
	delimiter         = []byte("/")
)

// BoltManager is the Manager for the Key-Value-Store Boltdb
type BoltManager struct {
	DB *bolt.DB
}

// newBoltManager is the factore method for BoltManager
func newBoltManager() *BoltManager {
	db, err := bolt.Open("gcr.db", 0600, nil)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "BoldManager"))
	}

	manager := &BoltManager{
		DB: db,
	}

	return manager
}

// PutRecord inserts or creates a new record
func (m *BoltManager) PutRecord(r *Record) error {
	oldRecord, err := m.FindNextRevision(r)
	if err != nil {
		return errors.Wrap(err, "BoltDB")
	}

	return m.InsertRevision(r, oldRecord)
}

// GetRecord a single entry
func (m *BoltManager) GetRecord(contentType, key []byte) (Record, error) {
	var content []byte
	m.DB.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(contentTypeBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}
		content = b.Get(key)
		return nil
	})

	return NewRecordFromJSON(content)
}

// DeleteRecord removes a single record
func (m *BoltManager) DeleteRecord(r Record) error {
	return m.DB.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(r.ContentType)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		// Delete from HEAD
		err = b.Delete(r.Key)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Delete record")
		}

		// Delete all revisions
		bh, err := b.CreateBucketIfNotExists(historyBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		c := bh.Cursor()
		prefix := m.generatePrefix(r.Key)
		for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			bh.Delete(k)
		}

		return nil
	})
}

// InsertRevision adds a new revision viewing the old key
func (m *BoltManager) InsertRevision(newRecord *Record, oldRecord Record) error {
	storeContent, err := newRecord.ToJSON()
	if err != nil {
		return errors.Wrap(err, "BoltDB")
	}

	err = m.DB.Batch(func(tx *bolt.Tx) error {
		// Get record bucket for contentType
		br, err := tx.CreateBucketIfNotExists(newRecord.ContentType)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		// Get history bucket for contentType
		bh, err := br.CreateBucketIfNotExists(historyBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		diffPatch, err := jsonpatch.CreateMergePatch(oldRecord.Content, newRecord.Content)
		if err != nil {
			return errors.Wrap(err, "BoltDB")
		}

		bh.Put(newRecord.TimedKey(), diffPatch)

		// Store/Update key
		return br.Put(newRecord.Key, storeContent)
	})

	return err
}

// FindNextRevision searches for the next revision (one older) of a given record
func (m *BoltManager) FindNextRevision(r *Record) (Record, error) {
	var value []byte
	m.DB.View(func(tx *bolt.Tx) error {
		bh := tx.Bucket(r.ContentType).Bucket(historyBucket)

		// Timed key is {key}/timestamp
		// This helps us because we can scan starting with key and go to the next key
		// Example: key = a timestamp = 900 Timedkey is "a/900"
		// More keys: a/600, a/300 and b/600
		// If we do a prefix scan we can find all "a/" prefixed and go to the next one
		// which has to be a/600 from a/300
		c := bh.Cursor()
		prefix := m.generatePrefix(r.Key)
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Prev() {
			fmt.Printf("FOUND: key=%s, value=%s\n", k, v)
			if k != nil {
				value = v
				return nil
			}
		}

		return nil
	})

	// No record found but is not an error!
	if len(value) == 0 {
		return Record{}, nil
	}

	return NewRecordFromJSON(value)
}

// PutContentType updates or creates a content type by key
func (m *BoltManager) PutContentType(c *ContentType) error {
	return m.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(contentTypeBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		return bucket.Put(c.Key, c.Validation)
	})
}

// GetContentType a single entry
func (m *BoltManager) GetContentType(key []byte) (ContentType, error) {
	var content []byte
	m.DB.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(contentTypeBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}
		content = b.Get(key)
		return nil
	})

	return NewContentTypeFromJSON(content)
}

// DeleteContentType removes a content type from storage and the bucket for contenttype
func (m *BoltManager) DeleteContentType(c ContentType) error {
	return m.DB.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(contentTypeBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		err = b.Delete(c.Key)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Delete content type")
		}

		err = b.DeleteBucket(c.Key)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Content type bucket delete")
		}

		return nil
	})
}

// generatePrefix generates the scan prefix for a given key and delimiter
func (m *BoltManager) generatePrefix(key []byte) []byte {
	return append(key[:], delimiter[:]...)
}
