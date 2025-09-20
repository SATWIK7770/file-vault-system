package storage

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
)

type LocalStorage struct {
    BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
    return &LocalStorage{BasePath: basePath}
}

func (s *LocalStorage) Save(fileHeader *multipart.FileHeader, destPath string) error {
    src, err := fileHeader.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    fullPath := filepath.Join(s.BasePath, destPath)
    if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
        return err
    }

    dst, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    _, err = io.Copy(dst, src)
    return err
}

func (s *LocalStorage) Get(path string) (string, error) {
    fullPath := filepath.Join(s.BasePath, path)
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        return "", fmt.Errorf("file not found")
    }
    return fullPath, nil
}
