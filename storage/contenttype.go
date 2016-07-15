package storage

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ContentType explains the content type object
type ContentType struct {
	Key        []byte
	Validation []byte
}

// NewContentType create a new content type
func NewContentType(key, validation []byte) ContentType {
	return ContentType{
		Key:        key,
		Validation: validation,
	}
}

// NewContentTypeFromJSON unmarshalls the content type json data and returns a ContentType
func NewContentTypeFromJSON(jsonBlob []byte) (ContentType, error) {
	var c ContentType
	err := json.Unmarshal(jsonBlob, &c)
	if err != nil {
		return c, errors.Wrap(err, "NewContentTypeFromJSON")
	}

	return c, nil
}

// ToJSON creates a json blob from content type
func (c ContentType) ToJSON() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return b, errors.Wrap(err, "ToJSON")
	}

	return b, nil
}
