package record

import (
	"time"

	"github.com/sbani/gcr/contenttype"
)

// Record represents the ready to use json from storage
type Record struct {
	Key         string    `storm:"id" json:"key" xml:"key" form:"key" valid:"required"`
	ContentType string    `storm:"index" json:"contentType" xml:"contentType" form:"contentType" valid:"required,alphanum"`
	Revision    time.Time `storm:"index" json:"revision" xml:"revision" form:"revision"`
	Content     string    `json:"content" xml:"content" form:"content" valid:"required,json"`
}

// NewRecord create a new record
func NewRecord(c contenttype.ContentType, key, content string) Record {
	return Record{
		Key:         key,
		ContentType: c.Key,
		Revision:    time.Now(),
		Content:     content,
	}
}
