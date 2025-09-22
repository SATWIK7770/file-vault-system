package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
	"errors"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IncrementExpectedStorage(userID uint, size int64) error {
    return r.db.Model(&models.User{}).
        Where("id = ?", userID).
        Update("expected_storage", gorm.Expr("expected_storage + ?", size)).Error
}

func (r *UserRepository) IncrementActualStorage(userID uint, size int64) error {
    return r.db.Model(&models.User{}).
        Where("id = ?", userID).
        Update("actual_storage", gorm.Expr("actual_storage + ?", size)).Error
}


func (r *UserRepository) GetUserStorageUsed(userID uint) (int64, error) {
    var user models.User
    if err := r.db.Select("actual_storage").First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return 0, errors.New("user not found")
        }
        return 0, err
    }
    return user.ActualStorage, nil
}


func (r *UserRepository) UpdateUserStorage(userID uint, actualDelta int64, expectedDelta int64) error {
    // Use GORM's Updates with expression to increment/decrement values
    return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
        "actual_storage":      gorm.Expr("actual_storage + ?", actualDelta),
        "expected_storage":  gorm.Expr("expected_storage + ?", expectedDelta),
    }).Error
}
