package db

import (
    "fmt"
    "log"

    migrate "github.com/golang-migrate/migrate"
    _ "github.com/golang-migrate/migrate/database/postgres"
    _ "github.com/golang-migrate/migrate/source/file"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/config"
)

func RunMigrations(cfg *config.Config) error {
    migrationPath := "file://migrations"
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode)

    m, err := migrate.New(migrationPath, dbURL)
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration failed: %w", err)
    }

    log.Println("Database migrated successfully!")
    return nil
}
