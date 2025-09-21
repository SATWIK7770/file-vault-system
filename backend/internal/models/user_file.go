package models

import(
	"time"
)

type UserFile struct {
    ID     uint `gorm:"primaryKey" json:"id"`

    UserID uint `gorm:"not null;uniqueIndex:user_file_idx" json:"user_id"`
    FileID uint `gorm:"not null;uniqueIndex:user_file_idx" json:"file_id"`
	FileName string  `gorm:"not null" json:"file_name"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
	DownloadTimes int  `gorm:"not null;default:0" json:"download_times"`

    // Relations with cascade delete
    User User `gorm:"constraint:OnDelete:CASCADE;" json:"user"`
    File File `gorm:"constraint:OnDelete:CASCADE;" json:"file"`
}

