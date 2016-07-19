package server

import (
	"github.com/labstack/echo"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/http"
	"github.com/sbani/gcr/storage"
)

// newContentTypeHandler bootraps the hnadler for the content type
func newContentTypeHandler(c *config.Config, e *echo.Echo, s storage.ContentTypeManager) *http.ContentTypeHandler {
	h := &http.ContentTypeHandler{
		Storage: s,
	}
	h.SetRoutes(e)

	return h
}
