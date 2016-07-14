package storage

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/evanphx/json-patch"
	"github.com/pkg/errors"
)

var (
	historyBucket     = []byte("history")
	contentTypeBucket = []byte("contenttype")
)

// BoltManager is the Manager for the Key-Value-Store Boltdb
type BoltManager struct {
	db *bolt.DB
}

// newBoltManager is the factore method for BoltManager
func newBoltManager() *BoltManager {
	db, err := bolt.Open("gcr.db", 0600, nil)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "BoldManager"))
	}

	manager := &BoltManager{
		db: db,
	}

	return manager
}

// RecordToJSON inserts a new record
func (m *BoltManager) RecordToJSON(contentType, content []byte) (Record, error) {
	newRecord := NewRecord(GenerateKey(), content)
	return m.InsertRevision(contentType, newRecord, Record{})
}

// RecordToJSONWithkey inserts a new record but adds a fixed key (e.g.)
func (m *BoltManager) RecordToJSONWithkey(contentType, key, content []byte) (Record, error) {
	newRecord := NewRecord(key, content)
	return m.InsertRevision(contentType, newRecord, Record{})
}

// InsertRevision adds a new revision viewing the old key
func (m *BoltManager) InsertRevision(contentType []byte, newRecord Record, oldRecord Record) (Record, error) {
	storeContent, err := RecordToJSON(newRecord)
	if err != nil {
		return newRecord, errors.Wrap(err, "BoltDB")
	}

	err = m.db.Batch(func(tx *bolt.Tx) error {
		// Get record bucket for contentType
		br, err := tx.CreateBucketIfNotExists(contentType)
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

		timedkey := fmt.Sprintf("%d/%s", time.Now().UnixNano(), newRecord.Key)
		bh.Put([]byte(timedkey), diffPatch)

		// Store/Update key
		return br.Put(newRecord.Key, storeContent)
	})

	return newRecord, err
}

// HasRecord returns if the record for key exists
func (m *BoltManager) HasRecord(contentType, key []byte) bool {
	var value []byte
	m.db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(contentType).Get(key)
		return nil
	})

	return value != nil
}

// CreateContentType adds a new content type
func (m *BoltManager) CreateContentType(key, content []byte) error {
	return m.UpdateContentType(key, content)
}

// UpdateContentType updates a content type by key
func (m *BoltManager) UpdateContentType(key, content []byte) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(contentTypeBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		return bucket.Put(key, content)
	})
}

// GetContentType a single entry
func (m *BoltManager) GetContentType(key []byte) ([]byte, error) {
	var content []byte
	err := m.db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(contentTypeBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}
		content = b.Get(key)
		return nil
	})

	return content, err
}

// DeleteContentType removes a content type from storage and the bucket for contenttype
func (m *BoltManager) DeleteContentType(key []byte) error {
	return m.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(contentTypeBucket)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Bucket")
		}

		err = b.Delete(key)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Delete content type")
		}

		err = b.DeleteBucket(key)
		if err != nil {
			return errors.Wrap(err, "BoltDB: Content type bucket delete")
		}

		return nil
	})
}
