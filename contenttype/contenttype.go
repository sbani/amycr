package contenttype

// ContentType explains the content type object
type ContentType struct {
	Key        string `storm:"id" json:"key" xml:"key" form:"key" valid:"required,alphanum"`
	Validation string `json:"validation" xml:"validation" form:"validation" valid:"required,json"`
}

// NewContentType create a new content type
func NewContentType(key, validation string) ContentType {
	return ContentType{
		Key:        key,
		Validation: validation,
	}
}
