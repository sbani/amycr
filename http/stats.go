package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/storage"
)

// StorageHandler is the Handler for contenttype
type StorageHandler struct {
	Storage storage.Manager
}

const (
	// StatisticHandlerPath is root path for all statistic actions
	StatsHandlerPath = "/stats"

	// BackupHandlerPath is the root path for all backup actions
	BackupHandlerPath = "/backup"
)

// SetRoutes adds the routes related to the handler
func (h *StorageHandler) SetRoutes(e *echo.Echo) {
	e.GET(StatsHandlerPath, h.All)
	e.GET(BackupHandlerPath, h.Download)
}

// Download api call to download the database data file
func (h *StorageHandler) Download(c echo.Context) error {
	err := h.Storage.BackupDownload(c.Response())
	if err != nil {
		return errors.Wrap(err, "Storage")
	}

	return nil
}

// All action shows all statistic data at once
func (h *StorageHandler) All(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Storage.GetStats())
}
