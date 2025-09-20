package models

import (
    "golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
    ID uint             `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"not null" json:"username"`
    Email     string    `gorm:"not null" json:"email"`
    Password  string    `gorm:"not null" json:"password"`
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
