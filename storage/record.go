package storage

import (
	"crypto/sha1"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Record represents the ready to use json from storage
type Record struct {
	Key      []byte
	Revision time.Time
	Content  []byte
}

// NewRecord create a new record
func NewRecord(Key []byte, content []byte) Record {
	return Record{
		Key:      Key,
		Revision: time.Now(),
		Content:  content,
	}
}

// RecordToJSON create json blog from record
func RecordToJSON(record Record) ([]byte, error) {
	b, err := json.Marshal(record)
	if err != nil {
		return b, errors.Wrap(err, "RecordToJSON")
	}

	return b, nil
}

// RecordFromJSON unmarshalls the record json data and returns a Record
func RecordFromJSON(jsonBlob []byte) (Record, error) {
	var record Record
	err := json.Unmarshal(jsonBlob, &record)
	if err != nil {
		return record, errors.Wrap(err, "RecordFromJSON")
	}

	return record, nil
}

var hasher = sha1.New()

// GenerateKey creates a unique Key for a record
func GenerateKey() []byte {
	return hasher.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
}
