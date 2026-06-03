package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/dahaa/internal/storage"
)

type ImageHandler struct {
	storage *storage.ImageStorage
}

func NewImageHandler(storage *storage.ImageStorage) *ImageHandler {
	return &ImageHandler{storage: storage}
}

func (h *ImageHandler) ServeImage(c *echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "filename is required")
	}
	path := h.storage.GetImagePath(filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}
	return c.File(path)
}

func (h *ImageHandler) UploadImage(c *echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no image file provided")
	}
	if err := h.storage.ValidateImage(fileHeader); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	filename, err := h.storage.SaveImage(fileHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save image")
	}
	return c.JSON(http.StatusOK, map[string]string{"filename": filename})
}
