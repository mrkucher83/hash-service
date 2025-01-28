package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"github.com/pressly/goose"
	"net/url"
	"os"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	Timeout  int
}

func NewPoolConfig(cfg *Config) (*pgxpool.Config, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(cfg.Username),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DbName,
		cfg.Timeout)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	return poolConfig, nil
}

func NewConnection(poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewDbInstance() {
	cfg := &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		Timeout:  5,
	}

	poolConfig, err := NewPoolConfig(cfg)
	if err != nil {
		logger.Fatal("Pool config error: %v\n", err)
	}
	poolConfig.MaxConns = 5

	c, err := NewConnection(poolConfig)
	if err != nil {
		logger.Fatal("Connect to database failed: %v\n", err)
	}
	logger.Info("Successful connection to the DB!")

	mdb, _ := sql.Open("postgres", poolConfig.ConnString())
	err = mdb.Ping()
	if err != nil {
		logger.Fatal("database migration ping error: %v\n", err)
	}

	err = goose.Up(mdb, "./internal/migrations")
	if err != nil {
		logger.Fatal("database migration up error: %v\n", err)
	}

	_, err = c.Exec(context.Background(), ";")
	if err != nil {
		logger.Fatal("database ping error: %v\n", err)
	}
}
