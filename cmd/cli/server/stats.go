package server

import (
	"github.com/labstack/echo"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/http"
	"github.com/sbani/gcr/storage"
)

// newContentTypeHandler bootraps the hnadler for the content type
func newStatsHandler(c *config.Config, e *echo.Echo, m storage.Manager) *http.StatsHandler {
	h := &http.StatsHandler{
		Storage: m,
	}
	h.SetRoutes(e)

	return h
}
