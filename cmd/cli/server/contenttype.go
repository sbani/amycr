package server

import (
	"github.com/labstack/echo"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/contenttype"
	"github.com/sbani/gcr/storage"
)

// newContentTypeHandler bootraps the hnadler for the content type
func newContentTypeHandler(c *config.Config, e *echo.Echo, m *storage.Manager) *contenttype.Handler {
	h := &contenttype.Handler{
		Storage: m,
	}
	h.SetRoutes(e)

	return h
}
