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
	historyBucket = []byte("history")
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
	newRecord := NewRecord(GenerateID(), content)
	return m.InsertRevision(contentType, newRecord, nil)
}

// RecordToJSONWithID inserts a new record but adds a fixed key (e.g.)
func (m *BoltManager) RecordToJSONWithID(contentType, id, content []byte) (Record, error) {
	newRecord := NewRecord(id, content)
	return m.InsertRevision(contentType, newRecord, nil)
}

// UpdateRecord updates record to a new version
func (m *BoltManager) UpdateRecord(contentType, id, content []byte) (Record, error) {

	return m.InsertRevision(contentType, id, oldId, content)
}

// InsertRevision adds a new revision viewing the old id
func (m *BoltManager) InsertRevision(contentType []byte, newRecord Record, oldRecord Record) (Record, error) {
	storeContent, err := RecordToJSON(newRecord)
	if err != nil {
		return errors.Wrap(err, "BoltDB Store")
	}

	err = m.db.Batch(func(tx *bolt.Tx) error {
		// Get record bucket for contentType
		br, err := tx.CreateBucketIfNotExists(contentType)
		if err != nil {
			return err
		}

		// Get history bucket for contentType
		bh, err := br.CreateBucketIfNotExists(historyBucket)
		if err != nil {
			return err
		}

		diffPatch, err := jsonpatch.CreateMergePatch(oldRecord.Content, newRecord.Content)
		if err != nil {
			return err
		}

		timedID := fmt.Sprintf("%d/%s", time.Now().UnixNano(), newRecord.ID)
		bh.Put([]byte(timedID), diffPatch)

		// Store/Update id
		return br.Put(newRecord.ID, storeContent)
	})

	return newRecord, err
}

// HasRecord returns if the record for id exists
func (m *BoltManager) HasRecord(contentType, id []byte) bool {
	var value []byte
	m.db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(contentType).Get(id)
		return nil
	})

	return value != nil
}
