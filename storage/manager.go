package storage

import "github.com/labstack/echo/engine"

// Manager interface for storage
type Manager interface {
	// Record holds the rocerd manager
	Record() RecordManager

	// ContentType holds the contenttype manager
	ContentType() ContentTypeManager

	// GetStats returns an interface which then will get json encoded for output
	GetStats() interface{}

	// BackupDownload writes to a given response
	BackupDownload(engine.Response) error
}
