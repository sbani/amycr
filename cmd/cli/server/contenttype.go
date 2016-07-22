package server

import (
	"github.com/labstack/echo"
	"github.com/sbani/amycr/config"
	"github.com/sbani/amycr/http"
	"github.com/sbani/amycr/storage"
)

// newContentTypeHandler bootraps the hnadler for the content type
func newContentTypeHandler(c *config.Config, e *echo.Echo, s storage.ContentTypeManager) *http.ContentTypeHandler {
	h := &http.ContentTypeHandler{
		Storage: s,
	}
	h.SetRoutes(e)

	return h
}
