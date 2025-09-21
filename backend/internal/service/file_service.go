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
		userFile, err := fs.fileRepo.CreateUserReference(userID, existingFile , fileHeader.Filename);
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
	userFile, err := fs.fileRepo.CreateUserReference(userID, newFile , fileHeader.Filename);
	if err != nil {
		// Clean up file and file record if UserFile creation fails
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


// Private helper methods

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
    userFile, err := fs.fileRepo.GetUserFileByID(userfileID, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("file not found")
        }
        return errors.New("failed to find user file relation")
    }
    fileID := userFile.FileID

    // Step 2: Delete user <-> file relation
    if err := fs.fileRepo.DeleteUserFile(userID, fileID); err != nil {
        return errors.New("failed to delete file relationship")
    }

    // Step 3: Check remaining references
    count, err := fs.fileRepo.CountFileReferences(fileID)
    if err != nil {
        return errors.New("failed to check file references")
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

