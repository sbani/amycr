package contenttype

import "github.com/sbani/gcr/storage"

// ContentType explains the content type form object
type ContentType struct {
	Key        string `json:"key" xml:"key" form:"key" valid:"required,alphanum"`
	Validation string `json:"validation" xml:"validation" form:"validation" valid:"required,json"`
}

// ToStorageContentType fills the storage content type with the
func (c *ContentType) ToStorageContentType() storage.ContentType {
	s := storage.NewContentType([]byte(c.Key), []byte(c.Validation))
	return s
}
