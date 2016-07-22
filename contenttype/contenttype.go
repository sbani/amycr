package contenttype

import "github.com/xeipuuv/gojsonschema"

// ContentType explains the content type object
type ContentType struct {
	Key        string               `storm:"id" json:"key" xml:"key" form:"key" valid:"required,alphanum"`
	Validation string               `json:"validation" xml:"validation" form:"validation" valid:"required,json"`
	Schema     *gojsonschema.Schema `json:"omitempty"`
}

// Cache stores ContentType in cache
func (c *ContentType) Cache() {
	Cache.Add(c)
}
