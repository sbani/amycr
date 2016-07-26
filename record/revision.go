package record

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Revision represents the ready to use json from storage
type Revision struct {
	ID        string    `storm:"id" json:"id" xml:"id" form:"id"`
	Key       string    `storm:"index" json:"key" xml:"key" form:"key"`
	CreatedAt time.Time `storm:"index" json:"createdAt" xml:"createdAt" form:"createdAt"`
	Diff      string    `json:"content" xml:"content" form:"content"`
}

// NewRevision create a new record revision
func NewRevision(new *Record, old *Record) (Revision, error) {
	var r Revision

	if new == nil {
		return r, errors.New("'new' is not allowed to be nil in NewRevision()")
	}

	var oldContent string
	if old != nil {
		oldContent = old.Content
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, new.Content, true)

	now := time.Now()

	r.ID = fmt.Sprintf("%s/%s", new.Key, now.String())
	r.Key = new.Key
	r.CreatedAt = now
	r.Diff = dmp.DiffToDelta(diffs)

	return r, nil
}
