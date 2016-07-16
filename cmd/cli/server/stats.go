package server

import (
	"github.com/labstack/echo"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/stats"
	"github.com/sbani/gcr/storage"
)

// newContentTypeHandler bootraps the hnadler for the content type
func newStatsHandler(c *config.Config, e *echo.Echo, m storage.Manager) *stats.Handler {
	h := &stats.Handler{
		Storage: m,
	}
	h.SetRoutes(e)

	return h
}
