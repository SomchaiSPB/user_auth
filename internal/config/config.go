package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	PostgresDBConfig
	SqliteDBConfig
	httpPort          string
	environment       string
	authJwtSecret     string
	storage           string
	withFakeData      bool
	withTableTruncate bool
}

func (c Config) WithTableTruncate() bool {
	return c.withTableTruncate
}

func (c Config) WithFakeData() bool {
	return c.withFakeData
}

type SqliteDBConfig struct {
	dbFile string
}

func (s SqliteDBConfig) DbFile() string {
	return s.dbFile
}

type PostgresDBConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func (p PostgresDBConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.host, p.port, p.user, p.password, p.dbName)
}

func (c Config) Storage() string {
	return c.storage
}

func (c Config) HttpPort() string {
	return c.httpPort
}

func (c Config) Environment() string {
	return c.environment
}

func (c Config) AuthJwtSecret() []byte {
	return []byte(c.authJwtSecret)
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") == "" {
			return nil, fmt.Errorf("env variables are empty: %w", err)
		}
	}

	postgresDb := PostgresDBConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbName:   os.Getenv("DB_NAME"),
	}

	sqliteDb := SqliteDBConfig{
		dbFile: os.Getenv("DB_SQLITE_FILE"),
	}

	withFakeData, err := strconv.ParseBool(os.Getenv("APP_WITH_FAKE_DATA"))

	if err != nil {
		withFakeData = false
	}

	withTruncate, err := strconv.ParseBool(os.Getenv("APP_WITH_TABLE_TRUNCATE"))

	if err != nil {
		withTruncate = false
	}

	return &Config{
		PostgresDBConfig:  postgresDb,
		SqliteDBConfig:    sqliteDb,
		httpPort:          os.Getenv("APP_HTTP_PORT"),
		environment:       os.Getenv("APP_ENV"),
		authJwtSecret:     os.Getenv("AUTH_JWT_SECRET"),
		storage:           os.Getenv("APP_STORAGE"),
		withFakeData:      withFakeData,
		withTableTruncate: withTruncate,
	}, nil
}
