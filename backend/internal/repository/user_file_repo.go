package repository

import (
    "backend/internal/models"  
    "gorm.io/gorm" 
	"errors"           
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

// Fetch user_file record where user is the owner
func (r *UserFileRepository) GetOwnerUserFile(userID, userfileID uint) (*models.UserFile, error) {
	var uf models.UserFile
	if err := r.db.Where("user_id = ? AND id = ? AND is_owner = TRUE", userID, userfileID).First(&uf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found or not owned by user")
		}
		return nil, err
	}
	return &uf, nil
}

// Update visibility and public token
func (r *UserFileRepository) UpdateVisibility(uf *models.UserFile, visibility string, token *string) error {
    uf.Visibility = visibility
    uf.PublicToken = token
    if err := r.db.Model(uf).Updates(map[string]interface{}{
        "visibility":   visibility,
        "public_token": token,
    }).Error; err != nil {
        return err
    }
    return nil
}


func (r *UserFileRepository) GetByPublicToken(token string) (*models.UserFile, error) {
    var uf models.UserFile
    err := r.db.Where("public_token = ? AND visibility = 'public'", token).First(&uf).Error
    if err != nil {
        return nil, err
    }
    return &uf, nil
}

func (r *UserFileRepository) IncrementDownloadTimes(uf *models.UserFile) error {
    return r.db.Model(uf).UpdateColumn("download_times", gorm.Expr("download_times + ?", 1)).Error
}

func (r *UserFileRepository) CountFileReferences(fileID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.UserFile{}).Where("file_id = ?", fileID).Count(&count).Error
    return count, err
}

func (r *UserFileRepository) GetUserFileByID(userfileID, userID uint) (*models.UserFile, error) {
    var uf models.UserFile
    if err := r.db.Where("id = ? AND user_id = ?", userfileID, userID).First(&uf).Error; err != nil {
        return nil, err
    }
    return &uf, nil
}


func (r *UserFileRepository) DeleteUserFile(userID, fileID uint) error {
    res := r.db.Where("user_id = ? AND file_id = ?", userID, fileID).Delete(&models.UserFile{})
    if res.Error != nil {
        return res.Error
    }
    if res.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}