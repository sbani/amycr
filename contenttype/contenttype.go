package contenttype

import (
	"fmt"

	"github.com/hashicorp/golang-lru"
	"github.com/xeipuuv/gojsonschema"
)

// ContentType explains the content type object
type ContentType struct {
	Key        string `storm:"id" json:"key" xml:"key" form:"key" valid:"required,alphanum"`
	Validation string `json:"validation" xml:"validation" form:"validation" valid:"required,json"`
}

const schemaCacheSize = 100

var schemaCache *lru.Cache

func init() {
	schemaCache, _ = lru.New(schemaCacheSize)
}

// Schema returns the prepared json schema
func (c *ContentType) Schema() *gojsonschema.Schema {
	if v, ok := schemaCache.Get(c.Key); ok {
		fmt.Println("Loaded schema FROM CACHE")
		return v.(*gojsonschema.Schema)
	}

	l := gojsonschema.NewStringLoader(c.Validation)
	schema, _ := gojsonschema.NewSchema(l)

	return schema
}

// SetSchema adds the schema to cache
func (c *ContentType) SetSchema(s *gojsonschema.Schema) {
	schemaCache.Add(c.Key, s)
}
