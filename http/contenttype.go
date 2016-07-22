package http

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sbani/amycr/contenttype"
	"github.com/sbani/amycr/storage"
)

// ContentTypeHandler is the Handler for contenttype
type ContentTypeHandler struct {
	Storage storage.ContentTypeManager
}

const (
	// ContentTypeHandlerPath is the root path for all contenttype actions
	ContentTypeHandlerPath = "/contenttype"
)

// SetRoutes adds the routes related to the handler
func (h *ContentTypeHandler) SetRoutes(e *echo.Echo) {
	e.POST(ContentTypeHandlerPath, h.Put)
	e.PUT(ContentTypeHandlerPath, h.Put)
	e.GET(ContentTypeHandlerPath, h.List)
	e.GET(ContentTypeHandlerPath+"/:key", h.Get)
	e.DELETE(ContentTypeHandlerPath+"/:key", h.Delete)
}

// Put api call to create or updates a content type.
// Expecting a post request
func (h *ContentTypeHandler) Put(c echo.Context) error {
	ct := new(contenttype.ContentType)

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
	if err := h.Storage.Put(ct); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusOK, ct)
}

// List api call lists all content types
func (h *ContentTypeHandler) List(c echo.Context) error {
	contentTypes, err := h.Storage.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusOK, contentTypes)
}

// Get api call to get content type info
func (h *ContentTypeHandler) Get(c echo.Context) error {
	ct, err := h.Storage.Get(c.Param("key"))
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			return ErrContentTypeNotFound
		default:
			return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
		}
	}

	return c.JSON(http.StatusOK, ct)
}

// Delete api action to delete a single content type
func (h *ContentTypeHandler) Delete(c echo.Context) error {
	ct, err := h.Storage.Get(c.Param("key"))
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			return c.JSON(http.StatusNotFound, "Content type not found")
		default:
			return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
		}
	}

	err = h.Storage.Delete(&ct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusNoContent, "")
}
