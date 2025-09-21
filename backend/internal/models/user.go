package models

import (
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID uint             `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"not null" json:"username"`
    Email     string    `gorm:"not null" json:"email"`
    Password  string    `gorm:"not null" json:"password"`
    actualStorage int64 `gorm:"not null , default:0" json:"actual_storage"` 
    expectedStorage int64 `gorm:"not null , default:0" json:"expected_storage"` 
}


// HashPassword hashes a plaintext password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash verifies a hashed password with its plaintext
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
