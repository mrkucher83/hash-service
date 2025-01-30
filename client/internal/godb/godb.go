package godb

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Instance struct {
	Db *pgxpool.Pool
}

type Record struct {
	Id   int
	Text string
}

type Resp struct {
	Id   int
	Hash string
}

func (i *Instance) Close() {
	i.Db.Close()
}

func (i *Instance) AddRecord(ctx context.Context, rec *Record) (*Resp, error) {
	query := `INSERT INTO hash_storage (hash_value) VALUES ($1) RETURNING id, hash_value;`

	var resp Resp
	err := i.Db.QueryRow(ctx, query, rec.Text).Scan(&resp.Id, &resp.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to insert record: %w", err)
	}

	return &resp, nil
}

func (i *Instance) GetRecordById(ctx context.Context, id int) (*Record, error) {
	query := `SELECT id, hash_value FROM hash_storage WHERE id = $1;`

	var rec Record
	err := i.Db.QueryRow(ctx, query, id).Scan(&rec.Id, &rec.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to get record by id: %w", err)
	}

	return &rec, nil
}
