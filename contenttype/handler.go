package contenttype

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/storage"
)

// Handler is the Handler for contenttype
type Handler struct {
	Storage *storage.Manager
}

const (
	// ContentTypeHandlerPath is the root path for all contenttype actions
	ContentTypeHandlerPath = "/contenttype"
)

// SetRoutes adds the routes related to the handler
func (h *Handler) SetRoutes(e *echo.Echo) {
	e.POST(ContentTypeHandlerPath, h.Put)
	e.PUT(ContentTypeHandlerPath, h.Put)
	e.GET(ContentTypeHandlerPath, h.Get)
}

// Put api call to create or updates a content type.
// Expecting a post request
func (h *Handler) Put(c echo.Context) error {
	ct := new(storage.ContentType)

	// Bind input
	if err := c.Bind(ct); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Wrap(err, "ContentType Put"))
	}

	// Validate input
	if v, err := govalidator.ValidateStruct(ct); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else if !v {
		return c.JSON(http.StatusBadRequest, errors.New("Payload did not validate."))
	}

	// Put to storage

	if err := h.Storage.PutContentType(ct); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Wrap(err, "Storage"))
	}

	return c.JSON(http.StatusOK, ct)
}

// Get api call to get content type info
func (h *Handler) Get(c echo.Context) error {
	callback := c.QueryParam("callback")
	var content struct {
		Response  string    `json:"response"`
		Timestamp time.Time `json:"timestamp"`
		Random    int       `json:"random"`
	}
	content.Response = "Sent via JSONP"
	content.Timestamp = time.Now().UTC()
	content.Random = rand.Intn(1000)
	return c.JSONP(http.StatusOK, callback, &content)
}
