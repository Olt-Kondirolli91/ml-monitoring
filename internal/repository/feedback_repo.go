package repository

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/Olt-Kondirolli91/ml-monitoring/internal/models"
)

type FeedbackRepository interface {
    InsertFeedback(ctx context.Context, fb models.Feedback) error
    GetFeedbackByInferenceID(ctx context.Context, inferenceID string) ([]models.Feedback, error)
}

type feedbackRepo struct {
    db *sql.DB
}

func NewFeedbackRepository(db *sql.DB) FeedbackRepository {
    return &feedbackRepo{db: db}
}


func (r *feedbackRepo) InsertFeedback(ctx context.Context, fb models.Feedback) error {
    query := `
        INSERT INTO feedback (id, inference_id, feedback_data)
        VALUES ($1, $2, $3::jsonb)
    `
    _, err := r.db.ExecContext(ctx, query, fb.ID, fb.InferenceID, fb.FeedbackData)
    return err
}

func (r *feedbackRepo) GetFeedbackByInferenceID(ctx context.Context, inferenceID string) ([]models.Feedback, error) {
    query := `
        SELECT id, inference_id, feedback_data, created_at
        FROM feedback
        WHERE inference_id = $1
    `
    rows, err := r.db.QueryContext(ctx, query, inferenceID)
    if err != nil {
        return nil, fmt.Errorf("GetFeedbackByInferenceID: %w", err)
    }
    defer rows.Close()

    var feedbacks []models.Feedback
    for rows.Next() {
        var fb models.Feedback
        if err := rows.Scan(&fb.ID, &fb.InferenceID, &fb.FeedbackData, &fb.CreatedAt); err != nil {
            return nil, err
        }
        feedbacks = append(feedbacks, fb)
    }
    return feedbacks, rows.Err()
}
