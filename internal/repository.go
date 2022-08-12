package internal

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

type imageInfo struct {
	name       string
	createdAt  time.Time
	modifiedAt time.Time
}

func (r *Repository) UploadImage(ctx context.Context, name string, data []byte, time time.Time) (string, error) {

	var id string
	query := "INSERT INTO images (name, raw, created_at, modified_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.DB.QueryRowContext(ctx, query, name, data, time, time).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) UpdateImage(ctx context.Context, id string, newName string, newData []byte, time time.Time) error {

	query := "UPDATE images SET (name, raw, modified_at) = ($1, $2, $3) WHERE id = $4"
	exec, err := r.DB.ExecContext(ctx, query, newName, newData, time, id)
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

func (r *Repository) DownloadImage(ctx context.Context, id string) ([]byte, error) {

	var data []byte
	query := "SELECT raw FROM images WHERE id = $1"
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) GetImagesList(ctx context.Context) ([]imageInfo, error) {

	var list []imageInfo
	query := "SELECT name, created_At, modified_at FROM images"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var imageInfo imageInfo
		err := rows.Scan(&imageInfo.name, &imageInfo.createdAt, &imageInfo.modifiedAt)
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
