package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/zizouhuweidi/dahaa/internal/storage"
)

// ImageHandler handles image-related HTTP requests
type ImageHandler struct {
	storage *storage.ImageStorage
}

// NewImageHandler creates a new image handler
func NewImageHandler(storage *storage.ImageStorage) *ImageHandler {
	return &ImageHandler{
		storage: storage,
	}
}

// ServeImage serves an image file
func (h *ImageHandler) ServeImage(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "filename is required")
	}

	// Get the full path to the image
	path := h.storage.GetImagePath(filename)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusNotFound, "image not found")
	}

	// Serve the file
	return c.File(path)
}

// UploadImage handles image uploads
func (h *ImageHandler) UploadImage(c echo.Context) error {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no image file provided")
	}

	// Validate the image
	if err := h.storage.ValidateImage(file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Save the image
	filename, err := h.storage.SaveImage(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save image")
	}

	// Return the filename
	return c.JSON(http.StatusOK, map[string]string{
		"filename": filename,
	})
}
