package tests

import (
    "context"
    "errors"
    "regexp"
    "testing"
	"time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/models"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/repository"
)

func TestInsertFeedback_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewFeedbackRepository(db)

    query := regexp.QuoteMeta(`INSERT INTO feedback (id, inference_id, feedback_data)
        VALUES ($1, $2, $3::jsonb)`)

    mock.ExpectExec(query).
        WithArgs("fb-id", "inf-id", `{"corrected":"output"}`).
        WillReturnResult(sqlmock.NewResult(1, 1))

    fb := models.Feedback{
        ID:           "fb-id",
        InferenceID:  "inf-id",
        FeedbackData: `{"corrected":"output"}`,
    }

    err = repo.InsertFeedback(context.Background(), fb)
    if err != nil {
        t.Errorf("InsertFeedback returned error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestInsertFeedback_FKError(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewFeedbackRepository(db)

    query := regexp.QuoteMeta(`INSERT INTO feedback (id, inference_id, feedback_data)
        VALUES ($1, $2, $3::jsonb)`)

    mock.ExpectExec(query).
        WithArgs("fb-id", "bad-inf-id", `{"test":"data"}`).
        WillReturnError(errors.New("foreign key constraint"))

    fb := models.Feedback{
        ID:           "fb-id",
        InferenceID:  "bad-inf-id",
        FeedbackData: `{"test":"data"}`,
    }

    err = repo.InsertFeedback(context.Background(), fb)
    if err == nil {
        t.Error("Expected foreign key error, got nil")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestGetFeedbackByInferenceID_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewFeedbackRepository(db)

    query := regexp.QuoteMeta(`SELECT id, inference_id, feedback_data, created_at
        FROM feedback
        WHERE inference_id = $1`)

    columns := []string{"id", "inference_id", "feedback_data", "created_at"}
    mock.ExpectQuery(query).
        WithArgs("inf-id").
        WillReturnRows(
            sqlmock.NewRows(columns).
                AddRow("fb-id-1", "inf-id", `{"corrected":"output1"}`, time.Now()).
                AddRow("fb-id-2", "inf-id", `{"corrected":"output2"}`, time.Now()),
        )

    feedbacks, err := repo.GetFeedbackByInferenceID(context.Background(), "inf-id")
    if err != nil {
        t.Errorf("GetFeedbackByInferenceID returned error: %v", err)
    }
    if len(feedbacks) != 2 {
        t.Errorf("Expected 2 feedback items, got %d", len(feedbacks))
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestGetFeedbackByInferenceID_Empty(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewFeedbackRepository(db)

    query := regexp.QuoteMeta(`SELECT id, inference_id, feedback_data, created_at
        FROM feedback
        WHERE inference_id = $1`)

    // Return no rows
    mock.ExpectQuery(query).
        WithArgs("inf-id").
        WillReturnRows(sqlmock.NewRows([]string{"id", "inference_id", "feedback_data", "created_at"}))

    feedbacks, err := repo.GetFeedbackByInferenceID(context.Background(), "inf-id")
    if err != nil {
        t.Errorf("Expected nil error, got %v", err)
    }
    if len(feedbacks) != 0 {
        t.Errorf("Expected 0 feedback items, got %d", len(feedbacks))
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}
