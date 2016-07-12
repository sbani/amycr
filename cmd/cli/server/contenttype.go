package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/contenttype"
	"github.com/sbani/gcr/storage"
)

// newContentTypeHandler bootraps the hnadler for the content type
func newContentTypeHandler(c *config.Config, router *httprouter.Router, m *storage.Manager) *contenttype.Handler {
	h := &contenttype.Handler{
		Manager: m
	}
	h.SetRoutes(router)
	return h
}
