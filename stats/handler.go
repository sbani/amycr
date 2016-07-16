package statistic

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sbani/gcr/storage"
)

// Handler is the Handler for contenttype
type Handler struct {
	Storage storage.Manager
}

const (
	// StatisticHandlerPath is the root path for all statistic actions
	StatisticHandlerPath = "/statistic"
)

// SetRoutes adds the routes related to the handler
func (h *Handler) SetRoutes(e *echo.Echo) {
	e.GET(StatisticHandlerPath, h.All)
}

// All action shows all statistic data at once
func (h *Handler) All(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Storage.GetStats())
}
