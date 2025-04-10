package db

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq" 
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/config"
)

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode,
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping db: %w", err)
    }
    return db, nil
}
