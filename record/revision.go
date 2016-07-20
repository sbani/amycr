package record

import (
	"time"

	"github.com/evanphx/json-patch"
	"github.com/pkg/errors"
)

// Revision represents the ready to use json from storage
type Revision struct {
	ID        int       `storm:"id" json:"id" xml:"id" form:"id"`
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

	var oldContent []byte
	if old != nil {
		oldContent = []byte(old.Content)
	}

	diffPatch, err := jsonpatch.CreateMergePatch(oldContent, []byte(new.Content))
	if err != nil {
		return r, errors.Wrap(err, "NewRevision")
	}

	r.Key = new.Key
	r.CreatedAt = time.Now()
	r.Diff = string(diffPatch)

	return r, nil
}
