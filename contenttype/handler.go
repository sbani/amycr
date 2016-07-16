package contenttype

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/storage"
)

// Handler is the Handler for contenttype
type Handler struct {
	Storage storage.Manager
}

const (
	// ContentTypeHandlerPath is the root path for all contenttype actions
	ContentTypeHandlerPath = "/contenttype"
)

// SetRoutes adds the routes related to the handler
func (h *Handler) SetRoutes(e *echo.Echo) {
	e.POST(ContentTypeHandlerPath, h.Put)
	e.PUT(ContentTypeHandlerPath, h.Put)
	e.GET(ContentTypeHandlerPath, h.List)
	e.GET(ContentTypeHandlerPath+"/:key", h.Get)
}

// Put api call to create or updates a content type.
// Expecting a post request
func (h *Handler) Put(c echo.Context) error {
	ct := new(ContentType)

	// Bind input
	if err := c.Bind(ct); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Validate input
	if v, err := govalidator.ValidateStruct(ct); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else if !v {
		return c.JSON(http.StatusBadRequest, errors.New("Payload did not validate.").Error())
	}

	// Put to storage
	sct := ct.ToStorageContentType()
	if err := h.Storage.PutContentType(&sct); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusOK, sct)
}

// List api call lists all content types
func (h *Handler) List(c echo.Context) error {
	contentTypes, err := h.Storage.ListContentTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusOK, contentTypes)
}

// Get api call to get content type info
func (h *Handler) Get(c echo.Context) error {
	key := []byte(c.Param("key"))

	ct, err := h.Storage.GetContentType(key)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
	}

	// empty key means
	if len(ct.Key) == 0 {
		return c.JSON(http.StatusNotFound, "Content type not found")
	}

	return c.JSON(http.StatusOK, ct)
}
