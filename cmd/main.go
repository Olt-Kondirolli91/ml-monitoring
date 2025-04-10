package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/Olt-Kondirolli91/ml-monitoring/internal/config"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/db"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/server"
)

func main() {
    // 1. Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    // 2. Run migrations
    if err := db.RunMigrations(cfg); err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }

    // 3. Connect to DB
    database, err := db.ConnectDB(cfg)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer database.Close()

    // 4. Create and start HTTP server
    srv := server.NewServer(database)
    go srv.Start("8080") // run in goroutine

    // 5. Shutdown handling
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Received shutdown signal")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server Shutdown Failed:%+v", err)
    }
    log.Println("Server exited properly")
}
