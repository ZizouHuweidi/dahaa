package httpapi

import (
	"net/http"
	"os"

	"github.com/zizouhuweidi/dahaa/internal/storage"
)

type ImageHandler struct {
	storage *storage.ImageStorage
}

func NewImageHandler(storage *storage.ImageStorage) *ImageHandler {
	return &ImageHandler{storage: storage}
}

func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	if filename == "" {
		writeError(w, http.StatusBadRequest, "filename is required")
		return
	}
	path := h.storage.GetImagePath(filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		writeError(w, http.StatusNotFound, "image not found")
		return
	}
	http.ServeFile(w, r, path)
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "no image file provided")
		return
	}
	_ = file.Close()

	fileHeader := r.MultipartForm.File["image"][0]
	if err := h.storage.ValidateImage(fileHeader); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	filename, err := h.storage.SaveImage(fileHeader)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save image")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"filename": filename})
}
