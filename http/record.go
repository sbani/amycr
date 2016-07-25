package http

import (
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sbani/amycr/record"
	"github.com/sbani/amycr/storage"
)

// RecordHandler is the Handler for record
type RecordHandler struct {
	Storage storage.Manager
}

const (
	// RecordHandlerPath is the root path for all record actions
	RecordHandlerPath = "/record"

	// RecordContentTypeHandlerPath is the root path for all record actions related to a content type
	RecordContentTypeHandlerPath = RecordHandlerPath + "/:contenttype/:key"
)

// SetRoutes adds the routes related to the handler
func (h *RecordHandler) SetRoutes(e *echo.Echo) {
	// Not content type based
	e.POST(RecordHandlerPath, h.Put)
	e.PUT(RecordHandlerPath, h.Put)

	// Group all routes which need content types and check (and set) the content type
	g := e.Group(RecordContentTypeHandlerPath)
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get content type
			ct, err := h.Storage.ContentType().Get(c.Param("contenttype"))
			if err != nil {
				switch err {
				case storm.ErrNotFound:
					return ErrContentTypeNotFound
				default:
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			}

			c.Set("contenttype", ct)

			// Get record
			r, err := h.Storage.Record().Get(ct.Key, c.Param("key"))
			if err != nil {
				switch err {
				case storm.ErrNotFound:
					return ErrRecordNotFound
				default:
					return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
				}
			}

			c.Set("record", r)

			return next(c)
		}
	})

	// Group related (content type based)
	g.GET("", h.Get)
	g.DELETE("", h.Delete)

	g.GET("/revisions", h.ListRevisions)
}

// Put api call to create or updates a content type.
// Expecting a post request
func (h *RecordHandler) Put(c echo.Context) error {
	var r record.Record
	r.Revision = time.Now()

	// Bind input
	if err := c.Bind(&r); err != nil {
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
	r := c.Get("record").(record.Record)
	return c.JSON(http.StatusOK, r)
}

// ListRevisions lists all revisions for a given record
func (h *RecordHandler) ListRevisions(c echo.Context) error {
	r := c.Get("record").(record.Record)
	revs, err := h.Storage.Record().GetRevisions(r)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.Wrap(err, "Storage").Error())
	}

	data := struct {
		Count        int               `json:"count"`
		HeadRevision time.Time         `json:"headRevision"`
		ContentType  string            `json:"contentType"`
		Key          string            `json:"key"`
		Revisions    []record.Revision `json:"revisions"`
	}{
		len(revs),
		r.Revision,
		r.ContentType,
		r.Key,
		revs,
	}

	return c.JSON(http.StatusOK, data)
}

// Delete api action to delete a single content type
func (h *RecordHandler) Delete(c echo.Context) error {
	r := c.Get("record").(record.Record)
	err := h.Storage.Record().Delete(r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "Storage").Error())
	}

	return c.JSON(http.StatusNoContent, "")
}
