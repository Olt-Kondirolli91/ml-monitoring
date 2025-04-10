package tests

import (
    "context"
    "errors"
    "sync"
    "time"

    "github.com/Olt-Kondirolli91/ml-monitoring/internal/models"
    "github.com/Olt-Kondirolli91/ml-monitoring/internal/repository"
)

var (
    // Some custom errors to simulate DB constraints in the mock
    errNotFound     = errors.New("not found")
    errForeignKey   = errors.New("foreign key constraint")
    errAlreadyExist = errors.New("record already exists")
)

// MockInferenceRepo is an in-memory implementation
type MockInferenceRepo struct {
    store map[string]models.Inference
    mu    sync.RWMutex
}

func NewMockInferenceRepo() repository.InferenceRepository {
    return &MockInferenceRepo{
        store: make(map[string]models.Inference),
    }
}

func (m *MockInferenceRepo) InsertInference(ctx context.Context, inf models.Inference) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // If it already exists, simulate a conflict
    if _, exists := m.store[inf.ID]; exists {
        return errAlreadyExist
    }
    inf.CreatedAt = time.Now()
    m.store[inf.ID] = inf
    return nil
}

func (m *MockInferenceRepo) UpdateHasFeedback(ctx context.Context, inferenceID string, hasFeedback bool) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    inf, ok := m.store[inferenceID]
    if !ok {
        return errNotFound
    }
    inf.HasFeedback = hasFeedback
    m.store[inferenceID] = inf
    return nil
}

func (m *MockInferenceRepo) GetInferenceByID(ctx context.Context, inferenceID string) (*models.Inference, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    inf, ok := m.store[inferenceID]
    if !ok {
        return nil, errNotFound
    }
    return &inf, nil
}

// MockFeedbackRepo is an in-memory implementation
type MockFeedbackRepo struct {
    store map[string][]models.Feedback
    mu    sync.RWMutex
}

func NewMockFeedbackRepo() repository.FeedbackRepository {
    return &MockFeedbackRepo{
        store: make(map[string][]models.Feedback),
    }
}

func (m *MockFeedbackRepo) InsertFeedback(ctx context.Context, fb models.Feedback) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Suppose if the "inference_id" is "bad-inf-id", simulate a foreign key failure
    if fb.InferenceID == "bad-inf-id" {
        return errForeignKey
    }

    fb.CreatedAt = time.Now()
    m.store[fb.InferenceID] = append(m.store[fb.InferenceID], fb)
    return nil
}

func (m *MockFeedbackRepo) GetFeedbackByInferenceID(ctx context.Context, inferenceID string) ([]models.Feedback, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    feedbacks, ok := m.store[inferenceID]
    if !ok {
        // Return empty slice if none
        return []models.Feedback{}, nil
    }
    return feedbacks, nil
}
