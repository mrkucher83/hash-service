package godb

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Instance struct {
	Db *pgxpool.Pool
}

type Record struct {
	Id     int
	Hashes []string
}

func (i *Instance) AddRecord(ctx context.Context, rec *Record) error {
	query := `INSERT INTO hash_storage (hash_value) VALUES ($1) RETURNING id;`

	err := i.Db.QueryRow(ctx, query, rec.Hashes).Scan(&rec.Id)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	return nil
}

func (i *Instance) GetRecordById(ctx context.Context, id int) (*Record, error) {
	query := `SELECT id, hash_value FROM hash_storage WHERE id = $1;`

	var rec Record
	err := i.Db.QueryRow(ctx, query, id).Scan(&rec.Id, &rec.Hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to get record by id: %w", err)
	}

	return &rec, nil
}
