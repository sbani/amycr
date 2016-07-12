package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/contenttype"
)

// Handler holds all other handlers and prepares them for routing
type Handler struct {
	ContentType *contenttype.Handler
}

// Start the handler and bootrap all others
func (h *Handler) Start(c *config.Config, router *httprouter.Router) {
	h.ContentType = newContentTypeHandler(c, router, storage)
}
