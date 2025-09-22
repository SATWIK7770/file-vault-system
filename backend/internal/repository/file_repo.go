package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
	"time"
	"errors"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Save creates a new file record in the database
func (r *FileRepository) Save(file *models.File) error {
	return r.db.Create(file).Error
}

// GetFilesByUser retrieves all files belonging to a specific user
func (r *FileRepository) GetFilesByUser(userID uint) ([]models.File, error) {
	var files []models.File
	if err := r.db.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// GetFileForDownload retrieves a file that belongs to a specific user
func (r *FileRepository) GetFileForDownload(userID uint, fileID uint) (*models.File, error) {
	var userFile models.UserFile

	// Step 1: Verify user-file relation exists
	if err := r.db.Where("user_id = ? AND file_id = ?", userID, fileID).First(&userFile).Error; err != nil {
		return nil, errors.New("file not found for this user")
	}

	if err := r.db.Model(&models.UserFile{}).Where("user_id = ? AND file_id = ?", userID, fileID).UpdateColumn("download_times", gorm.Expr("download_times + 1")).Error; err != nil {
    	return nil, errors.New("failed to update download count")
	}

	// Step 3: Fetch actual file info
	var file models.File
	if err := r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		return nil, errors.New("file metadata not found")
	}

	return &file, nil
}


// GetFileByHash finds a file by its hash (for deduplication)
func (r *FileRepository) GetFileByHash(hash string) (*models.File, error) {
	var file models.File
	err := r.db.Where("hash = ?", hash).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// UserHasFile checks if a user already has a file with the given hash
func (r *FileRepository) UserHasFile(userID uint, hash string) (bool, error) {
	var count int64
	err := r.db.Model(&models.File{}).
		Where("user_id = ? AND hash = ?", userID, hash).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateUserReference creates a UserFile relationship for deduplication
func (r *FileRepository) CreateUserReference(userID uint, existingFile *models.File, filename  string, ownership bool) (*models.UserFile, error) {
	// Create UserFile relationship
	userFile := &models.UserFile{
		UserID: userID,
		FileID: existingFile.ID,
		FileName : filename,
		UploadedAt : time.Now(),
		DownloadTimes : 0,
		IsOwner : ownership,
	}


	// Start transaction
	tx := r.db.Begin()

	// Create the user-file relationship
	if err := tx.Create(userFile).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Increment refCount
	if err := tx.Model(&models.File{}).Where("id = ?", existingFile.ID).
		Update("ref_count", gorm.Expr("ref_count + 1")).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return userFile, nil
}

// GetFileByID retrieves a file by its ID (without user restriction)
func (r *FileRepository) GetFileByID(fileID uint) (*models.File, error) {
	var file models.File
	err := r.db.First(&file, fileID).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}


func (r *FileRepository) DeleteFileRecord(fileID uint) error {
    return r.db.Where("id = ?", fileID).Delete(&models.File{}).Error
}

func (r *FileRepository) UpdateReferenceCount(fileID uint, count int64) error {
    return r.db.Model(&models.File{}).Where("id = ?", fileID).Update("ref_count", count).Error
}



// // GetFileStats returns statistics about files in the system
// func (r *FileRepository) GetFileStats() (map[string]interface{}, error) {
// 	stats := make(map[string]interface{})
	
// 	// Total file records
// 	var totalFileRecords int64
// 	if err := r.db.Model(&models.File{}).Count(&totalFileRecords).Error; err != nil {
// 		return nil, err
// 	}
// 	stats["total_file_records"] = totalFileRecords
	
// 	// Total user-file relationships
// 	var totalUserFiles int64
// 	if err := r.db.Model(&models.UserFile{}).Count(&totalUserFiles).Error; err != nil {
// 		return nil, err
// 	}
// 	stats["total_user_files"] = totalUserFiles
	
// 	// Total physical storage used
// 	var totalStorage int64
// 	if err := r.db.Model(&models.File{}).
// 		Select("COALESCE(SUM(size), 0)").
// 		Scan(&totalStorage).Error; err != nil {
// 		return nil, err
// 	}
// 	stats["total_physical_storage"] = totalStorage
	
// 	// Average reference count
// 	var avgRefCount float64
// 	if err := r.db.Model(&models.File{}).
// 		Select("COALESCE(AVG(ref_count), 0)").
// 		Scan(&avgRefCount).Error; err != nil {
// 		return nil, err
// 	}
// 	stats["avg_reference_count"] = avgRefCount
	
// 	return stats, nil
// }