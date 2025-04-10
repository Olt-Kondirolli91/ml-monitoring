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

func TestInsertInference_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewInferenceRepository(db)

    // The query your InsertInference method executes:
    query := regexp.QuoteMeta(`INSERT INTO inferences (id, model_name, model_version, input_data, output_data, has_feedback)
        VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6)`)

    mock.ExpectExec(query).
        WithArgs(
            "some-uuid",
            "test-model",
            "v1",
            `{"sample":"input"}`,
            `{"prediction":"output"}`,
            false,
        ).
        WillReturnResult(sqlmock.NewResult(1, 1))

    inf := models.Inference{
        ID:           "some-uuid",
        ModelName:    "test-model",
        ModelVersion: "v1",
        InputData:    `{"sample":"input"}`,
        OutputData:   `{"prediction":"output"}`,
        HasFeedback:  false,
    }

    err = repo.InsertInference(context.Background(), inf)
    if err != nil {
        t.Errorf("InsertInference returned error: %v", err)
    }

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestInsertInference_DBError(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewInferenceRepository(db)
    query := regexp.QuoteMeta(`INSERT INTO inferences (id, model_name, model_version, input_data, output_data, has_feedback)
        VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6)`)

    // Simulate a DB error
    mock.ExpectExec(query).
        WillReturnError(errors.New("DB failure"))

    inf := models.Inference{
        ID:           "some-uuid",
        ModelName:    "test-model",
        ModelVersion: "v1",
        InputData:    `{}`,
        OutputData:   `{}`,
        HasFeedback:  false,
    }

    err = repo.InsertInference(context.Background(), inf)
    if err == nil {
        t.Error("Expected DB error, got nil")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestGetInferenceByID_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewInferenceRepository(db)

    query := regexp.QuoteMeta(`SELECT id, model_name, model_version, input_data, output_data, created_at, has_feedback
        FROM inferences
        WHERE id = $1`)

    columns := []string{"id", "model_name", "model_version", "input_data", "output_data", "created_at", "has_feedback"}
    mock.ExpectQuery(query).
        WithArgs("some-inf-id").
        WillReturnRows(
            sqlmock.NewRows(columns).AddRow(
                "some-inf-id",
                "test-model",
                "v1",
                `{"sample":"input"}`,
                `{"prediction":"output"}`,
                time.Date(2025, 4, 10, 12, 0, 0, 0, time.UTC),
                false,
            ),
        )

    inf, err := repo.GetInferenceByID(context.Background(), "some-inf-id")
    if err != nil {
        t.Errorf("GetInferenceByID error: %v", err)
    }
    if inf == nil || inf.ID != "some-inf-id" {
        t.Errorf("Expected inference with ID 'some-inf-id', got %v", inf)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestGetInferenceByID_NotFound(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewInferenceRepository(db)

    query := regexp.QuoteMeta(`SELECT id, model_name, model_version, input_data, output_data, created_at, has_feedback
        FROM inferences
        WHERE id = $1`)

    // Return no rows
    mock.ExpectQuery(query).
        WithArgs("non-existent-id").
        WillReturnRows(sqlmock.NewRows([]string{
            "id", "model_name", "model_version", "input_data", "output_data", "created_at", "has_feedback",
        }))

    inf, err := repo.GetInferenceByID(context.Background(), "non-existent-id")
    if inf != nil {
        t.Error("Expected nil inference for non-existent ID")
    }
    if err == nil {
        t.Error("Expected an error for non-existent inference, got nil")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestUpdateHasFeedback_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewInferenceRepository(db)

    query := regexp.QuoteMeta(`UPDATE inferences 
        SET has_feedback = $1
        WHERE id = $2`)

    // Pretend 1 row is updated
    mock.ExpectExec(query).
        WithArgs(true, "some-inf-id").
        WillReturnResult(sqlmock.NewResult(0, 1)) // RowsAffected=1

    err = repo.UpdateHasFeedback(context.Background(), "some-inf-id", true)
    if err != nil {
        t.Errorf("UpdateHasFeedback returned error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}

func TestUpdateHasFeedback_NoRows(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to open sqlmock: %v", err)
    }
    defer db.Close()

    repo := repository.NewInferenceRepository(db)

    query := regexp.QuoteMeta(`UPDATE inferences 
        SET has_feedback = $1
        WHERE id = $2`)

    // 0 rows updated
    mock.ExpectExec(query).
        WithArgs(true, "non-existent-id").
        WillReturnResult(sqlmock.NewResult(0, 0))

    err = repo.UpdateHasFeedback(context.Background(), "non-existent-id", true)
    // If your code doesn't check RowsAffected(), you won't get an error.
    // Typically, you'd do something like:
    //  rows, _ := res.RowsAffected()
    //  if rows == 0 { return errors.New("no rows updated") }
    //
    // If you DO handle that, you'd expect an error. Let's assume we do:
    if err == nil {
        t.Error("Expected an error for 0 rows updated, got nil")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }
}
