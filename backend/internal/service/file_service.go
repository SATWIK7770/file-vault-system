package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	// "errors"
	"time"
    // "os"     
	"mime/multipart"    

	// "gorm.io/gorm"
)

type FileService struct {
	fileRepo *repository.FileRepository
}

func NewFileService(fileRepo *repository.FileRepository) *FileService {
	return &FileService{fileRepo: fileRepo}
}

// SaveFile stores file metadata in DB
func (fs *FileService) SaveFile(userID uint, fileHeader *multipart.FileHeader, storagePath string, hash string, UploadedAt time.Time) (*models.File, error) {
	file := &models.File{
		UserID:     userID,
		Filename:   fileHeader.Filename,
		Hash:       hash,
		StoragePath: storagePath,
		UploadedAt:  UploadedAt,
		Size:       fileHeader.Size, 
	}
	if err := fs.fileRepo.Save(file); err != nil {
		return nil, err
	}
	return file, nil
}

// List files by user
func (s *FileService) GetFilesByUser(userID uint) ([]models.File, error) {
    return s.fileRepo.GetFilesByUser(userID)
}

// Get a file by ID & owner
// func (fs *FileService) GetFile(id uint, userID uint) (*models.File, error) {
// 	var file models.File
// 	if err := fs.fileRepo.Where("id = ? AND user_id = ?", id, userID).First(&file).Error; err != nil {
// 		return nil, err
// 	}
// 	return &file, nil
// }

// Delete file metadata (and optionally remove from disk)
// func (fs *FileService) DeleteFile(id uint, userID uint) error {
// 	var file models.File
// 	if err := fs.fileRepo.Where("id = ? AND user_id = ?", id, userID).First(&file).Error; err != nil {
// 		return err
// 	}

// 	// Delete DB record
// 	if err := fs.fileRepo.Delete(&file).Error; err != nil {
// 		return err
// 	}

// 	// Remove file from disk
// 	if err := os.Remove(file.StoragePath); err != nil && !errors.Is(err, os.ErrNotExist) {
// 		return err
// 	}

// 	return nil
// }
