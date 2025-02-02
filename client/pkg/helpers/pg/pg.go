package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/mrkucher83/hash-service/client/internal/godb"
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

func NewConnection(ctx context.Context, poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func runMigrations(conf *pgxpool.Config) error {
	mdb, _ := sql.Open("postgres", conf.ConnString())
	err := mdb.Ping()
	if err != nil {
		return err
	}
	err = goose.Up(mdb, "/var")
	if err != nil {
		return err
	}
	return nil
}

func loadConfig() *Config {
	return &Config{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		DbName:   os.Getenv("DB_NAME"),
		Timeout:  5,
	}
}

func NewDbInstance() (*godb.Instance, error) {
	// Инициализация конфигурации
	cfg := loadConfig()

	// Настройка пула соединений
	poolConfig, err := NewPoolConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Pool config error: %v\n", err)
	}
	poolConfig.MaxConns = 5

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Установка соединения с базой данных
	c, err := NewConnection(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to database failed: %v\n", err)
	}
	logger.Info("Successful connection to the DB!")

	// Запуск миграций
	if err = runMigrations(poolConfig); err != nil {
		return nil, fmt.Errorf("database migration error: %v\n", err)
	}

	// Проверка подключения
	_, err = c.Exec(ctx, ";")
	if err != nil {
		return nil, fmt.Errorf("database ping error: %v\n", err)
	}

	return &godb.Instance{Db: c}, nil
}
