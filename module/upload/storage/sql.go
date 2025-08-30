package uploadstorage

import "gorm.io/gorm"

type uploadStore struct {
	db *gorm.DB
}

func NewUploadStore(db *gorm.DB) *uploadStore {
	return &uploadStore{db: db}
}
