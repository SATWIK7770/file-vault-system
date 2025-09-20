package models

import "time"

type File struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Filename    string    `gorm:"not null" json:"filename"`
	Size        int64     `gorm:"not null" json:"size"`
	Hash        string    `gorm:"not null" json:"hash"`
	StoragePath string    `gorm:"not null" json:"storage_path"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}
