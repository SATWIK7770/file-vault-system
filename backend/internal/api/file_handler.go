package api

import (
	"backend/internal/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler creates a new file handler
// Handler should only depend on service, not repository directly
func NewFileHandler(fs *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fs,
	}
}

// Upload handles file upload requests
// Responsibility: HTTP request/response handling only
func (h *FileHandler) Upload(c *gin.Context) {
	// 1. Extract user ID from context
	userID := c.GetUint("userID")
	fmt.Printf("handler: File:%d", userID)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. Get uploaded file from request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	// 3. Delegate all business logic to service layer
	result, err := h.fileService.ProcessFileUpload(userID, fileHeader)
	if err != nil {
<<<<<<< HEAD
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

	savedFile, err := h.fileService.SaveFile(userID, file, storagePath, hash, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db save failed"})
		return
	}

	if err := c.SaveUploadedFile(file, storagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusOK, savedFile)
=======
		// Handle different types of service errors with proper status codes
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

	// 4. Return success response with proper structure
   c.JSON(http.StatusOK, gin.H{
    "id":       result.ID,       // user_files.id
    "file_id":  result.FileID,   // files.id
    "filename": result.FileName, // user_files.file_name
	})
>>>>>>> aee3e54 (implemented deduplication of files , mime validation , deletion of files)
}

// ListFiles handles file listing requests
func (h *FileHandler) ListFiles(c *gin.Context) {
    userID := c.GetUint("userID")

<<<<<<< HEAD
	files, err := h.fileService.GetFilesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// GET /api/files/:id/download
func (h *FileHandler) DownloadFile(c *gin.Context) {

	userID := c.GetUint("userID")
	fileIDstr := c.Param("id")

    fileID64, err := strconv.ParseUint(fileIDstr, 10, 64)
=======
    files, err := h.fileService.GetFilesByUser(userID)
>>>>>>> aee3e54 (implemented deduplication of files , mime validation , deletion of files)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
        return
    }

<<<<<<< HEAD
	file, err := h.fileRepo.GetFileForDownload(userID , fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(file.Filename))
	c.Header("Content-Type", "application/octet-stream")

	c.File(file.StoragePath)
=======
    // Return only id + filename for frontend
    simplified := make([]gin.H, 0, len(files))
    for _, f := range files {
        simplified = append(simplified, gin.H{
        "id":        f.ID,        // UserFile row id (unique per user)
        "file_id":   f.FileID,    // reference to File
        "filename":  f.FileName,
        })
    }

    c.JSON(http.StatusOK, gin.H{"files": simplified})
}

// DownloadFile handles file download requests
func (h *FileHandler) DownloadFile(c *gin.Context) {
	// 1. Extract parameters
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

	// 2. Get file info from service
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

	// 3. Set appropriate headers and serve file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=\""+fileInfo.Filename+"\"")
	c.Header("Content-Type", fileInfo.MimeType)
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
    c.Header("Pragma", "no-cache")
    c.Header("Expires", "0")
	c.File(fileInfo.StoragePath)
>>>>>>> aee3e54 (implemented deduplication of files , mime validation , deletion of files)
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