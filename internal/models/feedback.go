package models

import "time"

type Feedback struct {
    ID          string    `json:"id"`
    InferenceID string    `json:"inference_id"`
    FeedbackData string   `json:"feedback_data"`
    CreatedAt   time.Time `json:"created_at"`
}
