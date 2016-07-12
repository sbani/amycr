package storage

import (
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

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

func (m *BoltManager) CreateRecord(contentType, content []byte) ([]byte, error) {
	return m.InsertRevision(contentType, GenerateID(), content)
}

func (m *BoltManager) CreateRecordWithId(contentType, id, content []byte) ([]byte, error) {
	return m.InsertRevision(contentType, id, content)
}

func (m *BoltManager) InsertRevision(contentType, id, content []byte) ([]byte, error) {
	err := m.db.Batch(func(tx *bolt.Tx) error {
		buf, err := CreateRecordJSON(id, content)

		// Add to record bucket
		br, err := tx.CreateBucketIfNotExists(contentType)
		if err != nil {
			return err
		}

		// Add to history bucket
		bh, err := tx.CreateBucketIfNotExists(contentType)
		buf, err := CreateRecordJSON(id, content)
		if err != nil {
			return err
		}

		return b.Put(id, buf)
	})

	return id, err
}

func (m *BoltManager) HasRecord(contentType, id []byte) bool {
	var value []byte
	m.db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(contentType).Get(id)
		return nil
	})

	return value != nil
}
