package storage

import (
	"github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/storage/boltdb"
)

// BoltManager is the Manager for the Key-Value-Store Boltdb
type BoltManager struct {
	ORM *storm.DB
	r   *boltdb.RecordStore
	c   *boltdb.ContentTypeStore
}

// NewBoltManager is the factore method for BoltManager
func NewBoltManager() *BoltManager {
	storm, err := storm.Open("gcr.db")
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "BoldManager"))
	}

	cs := boltdb.NewContentTypeStore(storm)

	manager := &BoltManager{
		ORM: storm,
		r:   boltdb.NewRecordStore(storm, cs),
		c:   cs,
	}

	return manager
}

// ContentType returns stats for the database
func (m *BoltManager) ContentType() ContentTypeManager {
	return m.c
}

// Record holds the rocerd manager
func (m *BoltManager) Record() RecordManager {
	return m.r
}

// GetStats returns stats for the database
func (m *BoltManager) GetStats() interface{} {
	return m.ORM.Bolt.Stats()
}
