package http

import (
	"net/http"

	"github.com/labstack/echo"
)

var (
	ErrContentTypeNotFound = echo.NewHTTPError(http.StatusNotFound, "ContentType not found")
	ErrRecordNotFound      = echo.NewHTTPError(http.StatusNotFound, "Record not found")
)
