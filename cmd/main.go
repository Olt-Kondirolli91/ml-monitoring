package main

import (
    "context"
    "log"

    "github.com/google/uuid"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/config"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/db"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/models"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/repository"
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

    // 4. Create repositories
    infRepo := repository.NewInferenceRepository(database)
    fbRepo := repository.NewFeedbackRepository(database)

    // 5. Insert a sample inference (for demonstration)
    infID := uuid.New().String() // generate a random ID
    inf := models.Inference{
        ID:           infID,
        ModelName:    "example_model",
        ModelVersion: "1.0.0",
        InputData:    `{"input":"some input data"}`,
        OutputData:   `{"output":"some output data"}`,
        HasFeedback:  false,
    }

    ctx := context.Background()

    err = infRepo.InsertInference(ctx, inf)
    if err != nil {
        log.Fatalf("Failed to insert inference: %v", err)
    }
    log.Printf("Inserted inference with ID: %s\n", infID)

    // 6. Insert sample feedback
    fb := models.Feedback{
        ID:          uuid.New().String(),
        InferenceID: infID,
        FeedbackData: `{"corrected_output":"the correct output"}`,
    }
    err = fbRepo.InsertFeedback(ctx, fb)
    if err != nil {
        log.Fatalf("Failed to insert feedback: %v", err)
    }
    log.Printf("Inserted feedback for inference ID: %s\n", infID)

    // 7. Update the inference has_feedback to true
    err = infRepo.UpdateHasFeedback(ctx, infID, true)
    if err != nil {
        log.Fatalf("Failed to update inference has_feedback: %v", err)
    }

    // 8. Fetch and print the inference + feedback
    insertedInf, err := infRepo.GetInferenceByID(ctx, infID)
    if err != nil {
        log.Fatalf("Failed to get inference: %v", err)
    }
    log.Printf("Fetched Inference: %+v\n", insertedInf)

    allFeedback, err := fbRepo.GetFeedbackByInferenceID(ctx, infID)
    if err != nil {
        log.Fatalf("Failed to get feedback: %v", err)
    }
    log.Printf("Feedback for Inference ID %s: %+v\n", infID, allFeedback)
}
