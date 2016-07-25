package boltdb

import (
	"github.com/asdine/storm"
	"github.com/pkg/errors"
	"github.com/sbani/amycr/record"
)

const (
	headBucketName     = "heads"
	revisionBucketName = "revisions"
)

// RecordStore is the store for records
type RecordStore struct {
	db *storm.DB
	cs *ContentTypeStore
}

// NewRecordStore creates a new recordstore
func NewRecordStore(db *storm.DB, cs *ContentTypeStore) *RecordStore {
	return &RecordStore{
		db: db,
		cs: cs,
	}
}

// Put inserts or creates a new contenttype
func (s *RecordStore) Put(rec record.Record) error {
	// Start a tx within the record's content type bucket
	tx, err := s.db.From(rec.ContentType).Begin(true)
	if err != nil {
		return err
	}

	head := tx.From(headBucketName)

	// Get current record
	current, err := s.Get(rec.ContentType, rec.Key)
	if err != nil && err != storm.ErrNotFound {
		return errors.Wrapf(err, "ContentType (%s)", rec.ContentType)
	}

	// Create revision
	rev, err := record.NewRevision(&rec, &current)
	if err != nil {
		return errors.Wrap(err, "Revision")
	}

	// Save revision
	revision := tx.From(revisionBucketName)
	err = revision.Save(&rev)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Revision")
	}

	// Update record
	err = head.Save(&rec)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Get return a single record entry
func (s *RecordStore) Get(contentType, key string) (record.Record, error) {
	var rec record.Record
	err := s.db.From(contentType).From(headBucketName).One("Key", key, &rec)

	return rec, err
}

// GetRevisions returns a list of revisions for a given record
func (s *RecordStore) GetRevisions(r record.Record) ([]record.Revision, error) {
	var revs []record.Revision
	err := s.db.From(r.ContentType).From(revisionBucketName).Find("Key", r.Key, &revs)

	return revs, err
}

// Delete removes a single record and all revesions
func (s *RecordStore) Delete(r record.Record) error {
	// Start tx
	tx, err := s.db.From(r.ContentType).Begin(true)
	if err != nil {
		return err
	}

	// Delete revisions
	revs, err := s.GetRevisions(r)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, revision := range revs {
		tx.From(revisionBucketName).Remove(&revision)
	}

	// Delete HEAD
	err = tx.From(headBucketName).Remove(&r)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
