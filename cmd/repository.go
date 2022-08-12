package cmd

import (
	"database/sql"
	"time"
)

type ImageInfo struct {
	Name       string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type Repository struct {
	DB *sql.DB
}

type RepositoryConfig struct {
	Host         string
	Port         string
	Username     string
	DatabaseName string
	Password     string
}

func NewRepository(dbSource string) (*Repository, error) {
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Repository{DB: db}, nil
}

func (r *Repository) UploadImage(name string, data []byte) (int64, error) {
	return 0, nil
}

func (r *Repository) UpdateImage(id int64, newName string, newData []byte) error {
	return nil
}

func (r *Repository) DownloadImage(id int64) ([]byte, error) {
	return nil, nil
}

func (r *Repository) GetImagesList() ([]ImageInfo, error) {
	return nil, nil
}
