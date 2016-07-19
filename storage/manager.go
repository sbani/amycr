package storage

// Manager interface for storage
type Manager interface {
	// Record holds the rocerd manager
	//Record() RecordManager

	// ContentType holds the contenttype manager
	ContentType() ContentTypeManager

	// GetStats returns an interface which then will get json encoded for output
	GetStats() interface{}
}
