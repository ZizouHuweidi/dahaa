package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type errorResponse struct {
	Error string `json:"error"`
}

type StrictJSONBinder struct{}

func (StrictJSONBinder) Bind(c *echo.Context, dst any) error {
	r := c.Request()
	if r.Body == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	return nil
}

func httpErrorHandler(c *echo.Context, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"
	var he *echo.HTTPError
	if errors.As(err, &he) {
		status = he.StatusCode()
		if he.Message != "" {
			message = he.Message
		} else {
			message = http.StatusText(status)
		}
	} else {
		slog.Error("request failed", "method", c.Request().Method, "path", c.Request().URL.Path, "error", err)
	}

	if err := c.JSON(status, errorResponse{Error: message}); err != nil {
		slog.Error("failed to write error response", "error", err)
	}
}
