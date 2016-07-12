package storage

import (
	"crypto/sha1"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// RecordJSON represents the ready to use json from storage
type RecordJSON struct {
	ID       []byte
	Revision time.Time
	Content  []byte
}

// CreateRecordJSON adds id and revision to the content given by user
func CreateRecordJSON(id []byte, content []byte) ([]byte, error) {
	var b []byte
	record := RecordJSON{
		ID:       id,
		Revision: time.Now(),
		Content:  content,
	}
	b, err := json.Marshal(record)
	if err != nil {
		return b, errors.Wrap(err, "CreateRecordJSON")
	}

	return b, nil
}

// GetRecordJSON unmarshalls the record json data
func GetRecordJSON(jsonBlob []byte) (RecordJSON, error) {
	var record RecordJSON
	err := json.UnMarshal(jsonBlob, &record)
	if err != nil {
		return record, errors.Wrap(err, "GetRecordJSON")
	}

	return record, nil
}

var hasher = sha1.New()

// GenerateID creates a unique id for a record
func GenerateID() []byte {
	return hasher.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
}
