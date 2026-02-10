package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found, using system environment variables")
	}
}

func getConnectionStr() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func NewDB() (*DB, error) {
	connStr := getConnectionStr()
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	slog.Info("Attempting to connect to database", "host", host, "database", dbname)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := conn.Ping(); err != nil {
		slog.Error("Database ping failed", "host", host, "error", err)
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	slog.Info("Successfully connected to database", "host", host)

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		slog.Info("Closing database connection")
		return db.conn.Close()
	}
	return nil
}

func (db *DB) GetConn() *sql.DB {
	return db.conn
}
