package contenttype

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sbani/gcr/storage"
)

// Handler is the Handler for contenttype
type Handler struct {
	Manager *storage.Manager
}

const (
	// ContentTypeHandlerPath is the root path for all contenttype actions
	ContentTypeHandlerPath = "/contenttype"
)

// SetRoutes adds the routes related to the handler
func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.POST(ContentTypeHandlerPath, h.Create)
}

// Create api call to create a content type.
// Expecting a post request
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

}
