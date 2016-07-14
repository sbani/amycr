package storage

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Record represents the ready to use json from storage
type Record struct {
	Key         []byte
	ContentType []byte
	Revision    time.Time
	Content     []byte
}

// NewRecord create a new record
func NewRecord(contentType, key, content []byte) Record {
	return Record{
		Key:         key,
		ContentType: contentType,
		Revision:    time.Now(),
		Content:     content,
	}
}

// NewRecordFromJSON unmarshalls the record json data and returns a Record
func NewRecordFromJSON(jsonBlob []byte) (Record, error) {
	var record Record
	err := json.Unmarshal(jsonBlob, &record)
	if err != nil {
		return record, errors.Wrap(err, "RecordFromJSON")
	}

	return record, nil
}

// ToJSON create json blog from record
func (r Record) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return b, errors.Wrap(err, "ToJSON")
	}

	return b, nil
}

// TimedKey creates a key with {key}/timestamp
func (r Record) TimedKey() []byte {
	return []byte(fmt.Sprintf("%s/%d", r.Key, r.Revision.UnixNano()))
}

var hasher = sha1.New()

// GenerateKey creates a unique Key for a record
func GenerateKey() []byte {
	return hasher.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
}
