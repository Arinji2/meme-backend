package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	pool *sql.DB
	once sync.Once
)

func initDB() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || database == "" {
		log.Fatal("Missing Database Environment Variables")
	}

	var err error
	pool, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database))
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	pool.SetConnMaxLifetime(time.Minute * 4)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		log.Println("Received interrupt signal. Closing database connection...")
		if err := pool.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		os.Exit(0)
	}()
}

func GetDB() *sql.DB {
	once.Do(initDB)
	return pool
}

func ExecuteQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	db := GetDB()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return rows, nil
}

func ExecuteQueryRow(ctx context.Context, query string, args ...interface{}) (*sql.Row, context.CancelFunc) {
	db := GetDB()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	return db.QueryRowContext(ctx, query, args...), cancel
}
