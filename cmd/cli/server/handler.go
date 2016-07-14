package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/contenttype"
	"github.com/sbani/gcr/pkg"
	"github.com/sbani/gcr/storage"
)

// Handler holds all other handlers and prepares them for routing
type Handler struct {
	ContentType *contenttype.Handler
}

// Start the handler and bootrap all others
func (h *Handler) Start(c *config.Config, router *httprouter.Router) {
	storage, err := storage.NewManager(c)
	if err != nil {
		pkg.Must(errors.Wrap(err, "Storage"))
	}

	h.ContentType = newContentTypeHandler(c, router, &storage)
}
