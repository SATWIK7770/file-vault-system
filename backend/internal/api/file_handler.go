package api

import (
	"backend/internal/service"
	"net/http"
	"strconv"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fs *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fs,
	}
}

// Upload handles file upload requests
func (h *FileHandler) Upload(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

    if err := h.fileService.CheckStorageQuota(userID, fileHeader.Size); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

	result, err := h.fileService.ProcessFileUpload(userID, fileHeader)
	if err != nil {
		switch err.Error() {
		case "mime mismatch":
			c.JSON(http.StatusBadRequest, gin.H{"error": "File type validation failed"})
		case "file already uploaded by user":
			c.JSON(http.StatusConflict, gin.H{"error": "You have already uploaded this file"})
		case "file too large":
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds limit"})
		case "file type not allowed":
			c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		case "failed to open uploaded file":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		case "failed to create file reference":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       result.ID,
		"file_id":  result.FileID,
		"filename": result.FileName,
	})
}

// ListFiles handles file listing requests
func (h *FileHandler) ListFiles(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    // Collect filters from query parameters
    filters := map[string]string{
        "filename":  c.Query("filename"),
        "mimeType":  c.Query("mimeType"),
        "uploader":  c.Query("uploader"),
        "minSize":   c.Query("minSize"),
        "maxSize":   c.Query("maxSize"),
        "startDate": c.Query("startDate"),
        "endDate":   c.Query("endDate"),
    }

    files, err := h.fileService.ListFilesForFrontend(userID.(uint), filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list files"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"files": files})
}



// DownloadFile handles file download requests
func (h *FileHandler) DownloadFile(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fileIDStr := c.Param("id")
	fileID64, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}
	fileID := uint(fileID64)

	fileInfo, err := h.fileService.GetFileForDownload(userID, fileID)
	if err != nil {
		switch err.Error() {
		case "file not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found or access denied"})
		case "file not found on disk":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Download failed"})
		}
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=\""+fileInfo.Filename+"\"")
	c.Header("Content-Type", fileInfo.MimeType)
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")	//avoid browser caching
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.File(fileInfo.StoragePath)
}

// DeleteFile handles file deletion requests
func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fileIDStr := c.Param("id")
	fileID64, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}
	userfileID := uint(fileID64)

	err = h.fileService.DeleteFile(userfileID, userID)
	if err != nil {
		switch err.Error() {
		case "file not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found or access denied"})
		case "failed to delete file relationship":
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete operation failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}


func (h *FileHandler) ChangeVisibility(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID64, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}
	userfileID := uint(fileID64)

	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	makePublicStr := c.Query("make_public")
	makePublic := makePublicStr == "true"

	log.Printf("Request body: make_public=%v", makePublic)

	uf, err := h.fileService.ChangeVisibility(userID, userfileID, makePublic)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	publicLink := ""
	if uf.Visibility == "public" && uf.PublicToken != nil {
		publicLink = "/download/" + *uf.PublicToken
	}

	response := gin.H{
        "file_id":     uf.ID,
        "visibility":  uf.Visibility,
        "public_link": publicLink,
    }
    log.Printf("Response: %+v", response)

	c.JSON(http.StatusOK, gin.H{
		"file_id":     uf.ID,
		"visibility":  uf.Visibility,
		"public_link": publicLink,
	})
}

func (h *FileHandler) DownloadPublic(c *gin.Context) {
    token := c.Param("token")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
        return
    }

    // Get both file metadata and user_file info
    file , err := h.fileService.GetFileByPublicToken(token)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

	fmt.Printf("DEBUG: Serving file from disk path: %s\n", file.StoragePath)

    // Serve file using user-specific filename
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=\""+file.Filename+"\"")
    c.Header("Content-Type", file.MimeType)
    c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
    c.Header("Pragma", "no-cache")
    c.Header("Expires", "0")

    c.File(file.StoragePath)
}

func (h *FileHandler) GetStorageStats(c *gin.Context) {
    userID := c.GetUint("userID")

    expected, actual, err := h.fileService.GetStorageStats(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user storage"})
        return
    }

    savings := expected - actual

    c.JSON(http.StatusOK, gin.H{
        "total_storage":      actual,
        "original_storage":   expected,
        "savings":            savings,
    })
}
