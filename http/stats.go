package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sbani/gcr/storage"
)

// StatsHandler is the Handler for contenttype
type StatsHandler struct {
	Storage storage.Manager
}

const (
	// StatisticHandlerPath is the root path for all statistic actions
	StatisticHandlerPath = "/stats"
)

// SetRoutes adds the routes related to the handler
func (h *StatsHandler) SetRoutes(e *echo.Echo) {
	e.GET(StatisticHandlerPath, h.All)
}

// All action shows all statistic data at once
func (h *StatsHandler) All(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Storage.GetStats())
}
