package storage

import "mime/multipart"

// Storage abstracts file saving/retrieving (local, S3, etc.)
type Storage interface {
    Save(fileHeader *multipart.FileHeader, destPath string) error
    Get(path string) (string, error) // return local path for download
}
