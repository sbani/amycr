package boltdb

import "github.com/asdine/storm"

// RecordStore is the store for records
type RecordStore struct {
	db *storm.DB
}

// NewRecordStore creates a new recordstore
func NewRecordStore(db *storm.DB) *RecordStore {
	return &RecordStore{
		db: db,
	}
}
