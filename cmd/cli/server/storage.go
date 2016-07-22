package server

import (
	"github.com/labstack/echo"
	"github.com/sbani/amycr/config"
	"github.com/sbani/amycr/http"
	"github.com/sbani/amycr/storage"
)

// newStorageHandler bootraps the hnadler for all storage actions
func newStorageHandler(c *config.Config, e *echo.Echo, m storage.Manager) *http.StorageHandler {
	h := &http.StorageHandler{
		Storage: m,
	}
	h.SetRoutes(e)

	return h
}
