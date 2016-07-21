package server

import (
	"github.com/labstack/echo"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/http"
	"github.com/sbani/gcr/storage"
)

// newStorageHandler bootraps the hnadler for all storage actions
func newStorageHandler(c *config.Config, e *echo.Echo, m storage.Manager) *http.StorageHandler {
	h := &http.StorageHandler{
		Storage: m,
	}
	h.SetRoutes(e)

	return h
}
