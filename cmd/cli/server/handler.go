package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/http"
	"github.com/sbani/gcr/pkg"
	"github.com/sbani/gcr/storage"
)

// Handler holds all other handlers and prepares them for routing
type Handler struct {
	e           *echo.Echo
	ContentType *http.ContentTypeHandler
	Stats       *http.StatsHandler
}

// Start the handler and bootrap all others
func (h *Handler) Start(c *config.Config, e *echo.Echo) {
	storage, err := storage.NewManager(c)
	if err != nil {
		pkg.Must(errors.Wrap(err, "Storage"))
	}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())

	h.ContentType = newContentTypeHandler(c, e, storage.ContentType())
	h.Stats = newStatsHandler(c, e, storage)
}
