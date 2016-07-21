package boltdb

import (
	"github.com/asdine/storm"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/record"
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
func (s *RecordStore) Put(rec *record.Record) error {
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
	rev, err := record.NewRevision(rec, &current)
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
	err = head.Save(rec)
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
func (s *RecordStore) GetRevisions(r *record.Record) ([]record.Revision, error) {
	var revs []record.Revision
	err := s.db.From(r.ContentType).From(revisionBucketName).Find("Key", r.Key, &revs)

	return revs, err
}

// Delete removes a single record and all revesions
func (s *RecordStore) Delete(r *record.Record) error {
	err := s.DeleteRevisions(r)
	if err != nil {
		return err
	}

	return s.db.From(r.ContentType).Remove(r)
}

// DeleteRevisions removes all revisions of a record
func (s *RecordStore) DeleteRevisions(r *record.Record) error {
	var revs []record.Revision
	err := s.db.From(r.ContentType).From(revisionBucketName).Find("Key", r.Key, &revs)
	if err != nil {
		// If somebody removed the revisions of that record before, it's ok to find anything
		if err == storm.ErrNotFound {
			return nil
		}

		return err
	}

	for _, revision := range revs {
		s.db.Remove(&revision)
	}

	return nil
}
