package api

import (
	"backend/internal/service"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	// "strconv"
	"time"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fs *service.FileService) *FileHandler {
	return &FileHandler{fileService: fs}
}

// POST /api/upload
func (h *FileHandler) Upload(c *gin.Context) {
	userID := c.GetUint("userID") // must match AuthMiddleware key
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	// build path
	storagePath := fmt.Sprintf("uploads/%s", file.Filename)

	// compute hash *before* saving metadata
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash file"})
		return
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	// save metadata + move file in service layer
	savedFile, err := h.fileService.SaveFile(userID, file, storagePath, hash, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db save failed"})
		return
	}

	// finally save to disk
	if err := c.SaveUploadedFile(file, storagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusOK, savedFile)
}


// internal/handler/file_handler.go
func (h *FileHandler) ListFiles(c *gin.Context) {
	userID := c.GetUint("userID")

	files, err := h.fileService.GetFilesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// GET /api/files/:id/download
func (h *FileHandler) DownloadFile(c *gin.Context) {
	// 1. Get file ID from URL
	id := c.Param("id")

	// 2. Find file in DB
	file, err := h.fileRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 3. Set headers so browser downloads instead of displaying
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(file.Filename))
	c.Header("Content-Type", "application/octet-stream")

	// 4. Serve the file
	c.File(file.StoragePath)
}

// DELETE /api/files/:id
// func (h *FileHandler) Delete(c *gin.Context) {
// 	userID := c.GetUint("userID")
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	err := h.fileService.DeleteFile(uint(id), userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"success": true})
// }
