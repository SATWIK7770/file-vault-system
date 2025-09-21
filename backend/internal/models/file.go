package models

type File struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Filename    string    `gorm:"not null" json:"filename"`
	Size        int64     `gorm:"not null" json:"size"`
	Hash        string    `gorm:"not null" json:"hash"`
	StoragePath string    `gorm:"not null" json:"storage_path"`
	MimeType 	string    `gorm:"not null" json:"mime_type"`
	RefCount   int64     `gorm:"not null , default:1" json:"ref_count"`
}
