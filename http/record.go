package http

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sbani/gcr/contenttype"
	"github.com/sbani/gcr/record"
	"github.com/sbani/gcr/storage"
)

// RecordHandler is the Handler for record
type RecordHandler struct {
	Storage storage.Manager
}

const (
	// RecordHandlerPath is the root path for all record actions
	RecordHandlerPath = "/record"

	// RecordContentTypeHandlerPath is the root path for all record actions related to a content type
	RecordContentTypeHandlerPath = RecordHandlerPath + "/:contenttype"
)

// SetRoutes adds the routes related to the handler
func (h *RecordHandler) SetRoutes(e *echo.Echo) {
	// Group all routes which need content types and check (and set) the content type
	g := e.Group(RecordContentTypeHandlerPath)
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				c.Error(err)
			}
			// Get content type
			ct, err := h.Storage.ContentType().Get(c.Param("contenttype"))
			if err != nil {
				switch err {
				case storm.ErrNotFound:
					return c.JSON(http.StatusNotFound, "Content type not found")
				default:
					return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
				}
			}

			fmt.Println(err, ct)

			c.Set("contenttype", ct)

			return err
		}
	})

	// Group related (content type based)
	g.GET("/:key", h.Get)
	g.DELETE("/:key", h.Delete)

	// Not content type based
	e.POST(RecordHandlerPath, h.Put)
	e.PUT(RecordHandlerPath, h.Put)
}

// Put api call to create or updates a content type.
// Expecting a post request
func (h *RecordHandler) Put(c echo.Context) error {
	r := new(record.Record)

	// Bind input
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Validate input
	if v, err := govalidator.ValidateStruct(r); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	} else if !v {
		return c.JSON(http.StatusBadRequest, errors.New("Payload did not validate.").Error())
	}

	// Put to storage
	if err := h.Storage.Record().Put(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusOK, r)
}

// Get api call to get content type info
func (h *RecordHandler) Get(c echo.Context) error {
	ct := c.Get("ct").(contenttype.ContentType)
	r, err := h.Storage.Record().Get(ct.Key, c.Param("key"))
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			return c.JSON(http.StatusNotFound, "Record not found")
		default:
			return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
		}
	}

	return c.JSON(http.StatusOK, r)
}

// Delete api action to delete a single content type
func (h *RecordHandler) Delete(c echo.Context) error {
	ct := c.Get("ct").(contenttype.ContentType)
	r, err := h.Storage.Record().Get(ct.Key, c.Param("key"))
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			return c.JSON(http.StatusNotFound, "Record not found")
		default:
			return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
		}
	}

	err = h.Storage.Record().Delete(&r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusNoContent, "")
}
