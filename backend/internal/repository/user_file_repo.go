package repository

import (
    "backend/internal/models"  
    "gorm.io/gorm" 
	"errors"      
	"time"  
	"strconv"   
)

type UserFileRepository struct {
	db *gorm.DB
}

func NewUserFileRepository(db *gorm.DB) *UserFileRepository {
	return &UserFileRepository{db: db}
}

type JoinedFileRow struct {
    UserFileID    uint
    FileID        uint
    FileName      string
    Size          int64
    MimeType      string
    UploaderName  string
    UploadedAt    time.Time
    Visibility    string
    IsOwner       bool
    DownloadTimes int
    PublicToken   *string
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

// user_file_repository.go
func (r *UserFileRepository) ListUserFilesWithFilters(userID uint, filters map[string]string) ([]JoinedFileRow, error) {
    var rows []JoinedFileRow
query := r.db.
    Table("user_files AS uf").
    Select(`uf.id AS user_file_id, 
            f.id AS file_id, 
            f.filename AS file_name, 
            f.size, 
            f.mime_type,
            u.username AS uploader_name,
            uf.uploaded_at, 
            uf.visibility, 
            uf.is_owner, 
            uf.download_times, 
            uf.public_token`).
    Joins("JOIN files f ON f.id = uf.file_id").
    Joins("JOIN users u ON u.id = uf.user_id").
    Where("uf.user_id = ?", userID)

    // Apply filters here as before
    if v, ok := filters["filename"]; ok && v != "" {
    query = query.Where("f.filename ILIKE ?", "%"+v+"%")
}

if v, ok := filters["mimeType"]; ok && v != "" {
    query = query.Where("f.mime_type ?", "%"+v+"%") 
}

if v, ok := filters["uploader"]; ok && v != "" {
    query = query.Where("owner.username ILIKE ?", "%"+v+"%") // match actual uploader
}

if v, ok := filters["minSize"]; ok && v != "" {
    if minSize, err := strconv.ParseInt(v, 10, 64); err == nil {
        query = query.Where("f.size >= ?", minSize)
    }
}

if v, ok := filters["maxSize"]; ok && v != "" {
    if maxSize, err := strconv.ParseInt(v, 10, 64); err == nil {
        query = query.Where("f.size <= ?", maxSize)
    }
}

if v, ok := filters["startDate"]; ok && v != "" {
    query = query.Where("uf.uploaded_at >= ?", v)
}

if v, ok := filters["endDate"]; ok && v != "" {
    query = query.Where("uf.uploaded_at <= ?", v)
}

if err := query.Scan(&rows).Error; err != nil {
    return nil, err
}

    return rows, nil
}
