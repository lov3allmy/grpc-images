package cmd

import (
	"database/sql"
	"errors"
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

func (r *Repository) UploadImage(name string, data []byte, time time.Time) (string, error) {

	var id string
	query := "INSERT INTO images (name, raw, created_at, modified_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.DB.QueryRow(query, name, data, time, time).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) UpdateImage(id string, newName string, newData []byte, time time.Time) error {

	query := "UPDATE images SET (name, raw, modified_at) = ($1, $2, $3) WHERE id = $4"
	exec, err := r.DB.Exec(query, newName, newData, time, id)
	if err != nil {
		return err
	}

	rows, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("nothing to update")
	}

	return nil
}

func (r *Repository) DownloadImage(id string) ([]byte, error) {

	var data []byte
	query := "SELECT raw FROM images WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) GetImagesList() ([]ImageInfo, error) {

	var list []ImageInfo
	query := "SELECT name, created_At, modified_at FROM images"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var imageInfo ImageInfo
		err := rows.Scan(&imageInfo.Name, &imageInfo.CreatedAt, &imageInfo.ModifiedAt)
		if err != nil {
			return list, err
		}
		list = append(list, imageInfo)
	}
	if err = rows.Err(); err != nil {
		return list, err
	}

	return list, nil
}
