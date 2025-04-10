package repository

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/Olt-Kondirolli91/ml-monitoring/internal/models"
)

type InferenceRepository interface {
    InsertInference(ctx context.Context, inf models.Inference) error
    UpdateHasFeedback(ctx context.Context, inferenceID string, hasFeedback bool) error
    GetInferenceByID(ctx context.Context, inferenceID string) (*models.Inference, error)
}

type inferenceRepo struct {
    db *sql.DB
}

func NewInferenceRepository(db *sql.DB) InferenceRepository {
    return &inferenceRepo{db: db}
}

func (r *inferenceRepo) InsertInference(ctx context.Context, inf models.Inference) error {
    query := `
        INSERT INTO inferences (id, model_name, model_version, input_data, output_data, has_feedback)
        VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6)
    `
    _, err := r.db.ExecContext(ctx, query,
        inf.ID, inf.ModelName, inf.ModelVersion, inf.InputData, inf.OutputData, inf.HasFeedback)
    return err
}

// UpdateHasFeedback updates the has_feedback flag for a given inference ID
func (r *inferenceRepo) UpdateHasFeedback(ctx context.Context, inferenceID string, hasFeedback bool) error {
    query := `
        UPDATE inferences 
        SET has_feedback = $1
        WHERE id = $2
    `
    res, err := r.db.ExecContext(ctx, query, hasFeedback, inferenceID)
    if err != nil {
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return fmt.Errorf("no rows updated for inference ID: %s", inferenceID)
    }
    return nil
}


func (r *inferenceRepo) GetInferenceByID(ctx context.Context, inferenceID string) (*models.Inference, error) {
    query := `
        SELECT id, model_name, model_version, input_data, output_data, created_at, has_feedback
        FROM inferences
        WHERE id = $1
    `
    row := r.db.QueryRowContext(ctx, query, inferenceID)
    var inf models.Inference
    err := row.Scan(&inf.ID, &inf.ModelName, &inf.ModelVersion, &inf.InputData,
        &inf.OutputData, &inf.CreatedAt, &inf.HasFeedback)
    if err != nil {
        return nil, fmt.Errorf("GetInferenceByID: %w", err)
    }
    return &inf, nil
}
