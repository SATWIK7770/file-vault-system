package repository

import (
    "backend/internal/models"  
    "gorm.io/gorm"            
)

type UserFileRepository struct {
	db *gorm.DB
}

func NewUserFileRepository(db *gorm.DB) *UserFileRepository {
	return &UserFileRepository{db: db}
}

func (r *UserFileRepository) UserHasFile(userID uint, hash string) (bool, error) {
    var count int64
    err := r.db.
        Model(&models.UserFile{}).
        Joins("JOIN files ON user_files.file_id = files.id").
        Where("user_files.user_id = ? AND files.hash = ?", userID, hash).
        Count(&count).Error

    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (r *UserFileRepository) GetUserFiles(userID uint) ([]models.UserFile, error) {
    var userFiles []models.UserFile
    if err := r.db.Where("user_id = ?", userID).Find(&userFiles).Error; err != nil {
        return nil, err
    }
    return userFiles, nil
}

