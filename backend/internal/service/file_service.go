package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"gorm.io/gorm"
	"crypto/rand"
)

type FileService struct {
    fileRepo     *repository.FileRepository
    userFileRepo *repository.UserFileRepository
    userRepo     *repository.UserRepository   
    config       FileConfig
}


type FileConfig struct {
	MaxFileSize   int64  // Maximum file size in bytes
	UploadDir     string // Directory to store uploaded files
	AllowedTypes  []string // Allowed MIME types (empty means all allowed)
}

func NewFileService(
    fileRepo *repository.FileRepository,
    userFileRepo *repository.UserFileRepository,
    userRepo *repository.UserRepository,
    config FileConfig,
) *FileService {
    return &FileService{
        fileRepo:     fileRepo,
        userFileRepo: userFileRepo,
        userRepo:     userRepo,
        config:       config,
    }
}


// ProcessFileUpload handles the complete file upload business logic
func (fs *FileService) ProcessFileUpload(userID uint, fileHeader *multipart.FileHeader) (*models.UserFile, error) {
	fmt.Printf("service: File:%d", userID)
	// 1. Validate file size
	if err := fs.validateFileSize(fileHeader.Size); err != nil {
		return nil, err
	}

	// 2. Open and validate file content
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.New("failed to open uploaded file")
	}
	defer file.Close()

	// 3. Validate MIME type
	mimeType, err := fs.validateMimeType(file, fileHeader.Filename)
	if err != nil {
		return nil, err
	}
	file.Seek(0, io.SeekStart) // Reset file pointer

	// 4. Compute file hash for deduplication
	hash, err := fs.computeFileHash(file)
	if err != nil {
		return nil, err
	}
	file.Seek(0, io.SeekStart) // Reset file pointer

	// 5. Check if user already uploaded this file
	exists, err := fs.userFileRepo.UserHasFile(userID, hash)
	if err != nil {
		return nil, errors.New("failed to check for existing file")
	}
	if exists {
		return nil, errors.New("file already uploaded by user")
	}

	// 6. Check for global file deduplication
	existingFile, err := fs.fileRepo.GetFileByHash(hash)
	if err == nil && existingFile != nil {
		// File exists globally, create user reference
		userFile, err := fs.fileRepo.CreateUserReference(userID, existingFile , fileHeader.Filename, false);
		if err != nil {
			return nil, errors.New("failed to create file reference")
		}
		
		_ = fs.userRepo.IncrementExpectedStorage(userID, existingFile.Size)
		// Return the existing file info (user will see it as their file)
		return userFile, nil
	}

	// 7. Generate unique storage path
	storagePath, err := fs.generateStoragePath(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	// 8. Save file to disk
	if err := fs.saveFileToDisk(file, storagePath); err != nil {
		return nil, err
	}

	// 9. Save file metadata to database with initial refCount of 1
	newFile := &models.File{ 
		Filename:    fileHeader.Filename,
		Hash:        hash,
		StoragePath: storagePath,
		MimeType:    mimeType,
		Size:        fileHeader.Size,
		RefCount:    0, // Initial reference count
	}

	if err := fs.fileRepo.Save(newFile); err != nil {
		// Clean up file if database save fails
		os.Remove(storagePath)
		return nil, errors.New("failed to save file metadata")
	}

	// Create the UserFile relationship
	userFile, err := fs.fileRepo.CreateUserReference(userID, newFile , fileHeader.Filename, true);
	if err != nil {
		os.Remove(storagePath)
		fs.fileRepo.DeleteFileRecord(newFile.ID) // Clean up the file record
		return nil, errors.New("failed to create user file relationship")
	}

	_ = fs.userRepo.IncrementExpectedStorage(userID, newFile.Size)
	_ = fs.userRepo.IncrementActualStorage(userID, newFile.Size)

	return userFile, nil
}

// GetFilesByUser retrieves all files for a user
func (fs *FileService) GetFilesByUser(userID uint) ([]models.UserFile, error) {
    return fs.userFileRepo.GetUserFiles(userID)

}


type FileFrontend struct {
	ID            uint   `json:"id"`
	FileID        uint   `json:"file_id"`
	Filename      string `json:"filename"`
	Size          int64  `json:"size"`
	Uploader      string `json:"uploader"`
	UploadDate    string `json:"upload_date"`
	IsPublic      string `json:"public"` // "yes"/"no"
	DownloadCount *int   `json:"downloads,omitempty"`
	Actions       struct {
		CanDownload     bool `json:"canDownload"`
		CanMakePublic   bool `json:"canMakePublic"`
		CanDelete       bool `json:"canDelete"`
		ShowDownloadCount bool `json:"showDownloadCount"`
	} `json:"actions"`
	PublicLink string `json:"public_link,omitempty"`
}

func (fs *FileService) ListFilesForFrontend(userID uint) ([]FileFrontend, error) {
	userFiles, err := fs.userFileRepo.GetUserFiles(userID)
	if err != nil {
		return nil, err
	}

	result := make([]FileFrontend, 0, len(userFiles))

	for _, uf := range userFiles {
		file, err := fs.fileRepo.GetFileByID(uf.FileID)
		if err != nil {
			continue
		}

		uploader := "exists in server"
		if uf.IsOwner {
			uploader = "you"
		}

		isPublic := "no"
		showDownloadCount := false
		publicLink := ""
		if uf.Visibility == "public" {
			isPublic = "yes"
			showDownloadCount = true
			if uf.PublicToken != nil {
				publicLink = fmt.Sprintf("/public/%s", *uf.PublicToken)
			}
		}

		var downloadCount *int
		if showDownloadCount {
			downloadCount = &uf.DownloadTimes
		}

		f := FileFrontend{
			ID:            uf.ID,
			FileID:        uf.FileID,
			Filename:      uf.FileName,
			Size:          file.Size,
			Uploader:      uploader,
			UploadDate:    uf.UploadedAt.Format("2006-01-02"),
			IsPublic:      isPublic,
			DownloadCount: downloadCount,
			PublicLink:    publicLink,
		}

		f.Actions.CanDownload = true
		f.Actions.CanMakePublic = uf.IsOwner
		f.Actions.CanDelete = true
		f.Actions.ShowDownloadCount = showDownloadCount

		result = append(result, f)
	}

	return result, nil
}


// GetFileForDownload retrieves file information for download
func (fs *FileService) GetFileForDownload(userID uint, fileID uint) (*models.File, error) {
	file, err := fs.fileRepo.GetFileForDownload(userID, fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}

	// Verify file exists on disk
	if _, err := os.Stat(file.StoragePath); os.IsNotExist(err) {
		return nil, errors.New("file not found on disk")
	}

	return file, nil
}


func (fs *FileService) validateFileSize(size int64) error {
	if fs.config.MaxFileSize > 0 && size > fs.config.MaxFileSize {
		return errors.New("file too large")
	}
	return nil
}

func (fs *FileService) validateMimeType(file multipart.File, filename string) (string, error) {
	// Detect MIME type from file content
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	detectedMime := http.DetectContentType(buf[:n])

	// Get expected MIME type from file extension
	expectedMime := mime.TypeByExtension(filepath.Ext(filename))

	// Validate if expected MIME type exists and matches
	if expectedMime != "" && detectedMime != expectedMime {
		return "", errors.New("mime mismatch")
	}

	// Check against allowed types if configured
	if len(fs.config.AllowedTypes) > 0 {
		allowed := false
		for _, allowedType := range fs.config.AllowedTypes {
			if detectedMime == allowedType {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", errors.New("file type not allowed")
		}
	}

	return detectedMime, nil
}

func (fs *FileService) computeFileHash(file multipart.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", errors.New("failed to compute file hash")
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (fs *FileService) generateStoragePath(filename string) (string, error) {
	// Ensure upload directory exists
	if err := os.MkdirAll(fs.config.UploadDir, 0755); err != nil {
		return "", errors.New("failed to create upload directory")
	}

	// Generate unique filename to avoid conflicts
	ext := filepath.Ext(filename)
	baseName := filename[:len(filename)-len(ext)]
	uniqueFilename := fmt.Sprintf("%s%s", baseName, ext)
	
	return filepath.Join(fs.config.UploadDir, uniqueFilename), nil
}

func (fs *FileService) saveFileToDisk(file multipart.File, storagePath string) error {
	dst, err := os.Create(storagePath)
	if err != nil {
		return errors.New("failed to create file on disk")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(storagePath) // Clean up on failure
		return errors.New("failed to save file to disk")
	}

	return nil
}

//
func (fs *FileService) DeleteFile(userfileID, userID uint) error {
    // Step 1: Fetch user_file entry to get file_id
    userFile, err := fs.userFileRepo.GetUserFileByID(userfileID, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("file not found")
        }
        return errors.New("failed to find user file relation")
    }
    fileID := userFile.FileID

    // Step 2: Delete user <-> file relation
    if err := fs.userFileRepo.DeleteUserFile(userID, fileID); err != nil {
        return errors.New("failed to delete file relationship")
    }

    // Step 3: Check remaining references
	count, err := fs.userFileRepo.CountFileReferences(fileID)
	if err != nil {
		return errors.New("failed to check file references")
	}

	// Step 3b: Update reference count in files table
	if err := fs.fileRepo.UpdateReferenceCount(fileID, count); err != nil {
		return errors.New("failed to update file reference count")
	}

    if count == 0 {
        // Step 4: Get file metadata
        file, err := fs.fileRepo.GetFileByID(fileID)
        if err != nil {
            return errors.New("file metadata not found")
        }

        // Step 5: Delete file record
        if err := fs.fileRepo.DeleteFileRecord(fileID); err != nil {
            return errors.New("failed to delete file record")
        }

        // Step 6: Delete file from disk
        if err := os.Remove(file.StoragePath); err != nil && !os.IsNotExist(err) {
            return fmt.Errorf("failed to delete file from disk: %w", err)
        }
    }

    return nil
}


func generatePublicToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Change file visibility; only owner allowed
func (fs *FileService) ChangeVisibility(userID, userfileID uint, makePublic bool) (*models.UserFile, error) {
	uf, err := fs.userFileRepo.GetOwnerUserFile(userID, userfileID)
	if err != nil {
		return nil, err
	}

	var token *string
	if makePublic {
		t := generatePublicToken()
		token = &t
	}

	visibility := "private"
	if makePublic {
		visibility = "public"
	}

	if err := fs.userFileRepo.UpdateVisibility(uf, visibility, token); err != nil {
		return nil, err
	}

	updated, err := fs.userFileRepo.GetOwnerUserFile(userID, userfileID)
	if err != nil {
    return nil, err
	}

	return updated, nil
}


// It validates that the token exists and the file is still public.
func (fs *FileService) GetFileByPublicToken(token string) (*models.File, error) {
    // Step 1: Lookup user_files entry via repository
    uf, err := fs.userFileRepo.GetByPublicToken(token)
    if err != nil {
        return nil , errors.New("link invalid or file private")
    }

    // Step 2: Increment download counter for analytics
    if err := fs.userFileRepo.IncrementDownloadTimes(uf); err != nil {
        // log error but still continue to serve the file
        fmt.Printf("Warning: failed to increment download count: %v\n", err)
    }

    // Step 3: Fetch actual file info from files table via repository
    file, err := fs.fileRepo.GetFileByID(uf.FileID)
    if err != nil {
        return nil , errors.New("file metadata not found")
    }

    return file, nil
}

// file_service.go
func (s *FileService) GetStorageStats(userID uint) (int64, int64, error) {
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return 0, 0, err
    }
    return user.ExpectedStorage, user.ActualStorage, nil
}
