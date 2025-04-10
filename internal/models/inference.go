package models

import "time"

type Inference struct {
    ID          string    `json:"id"`
    ModelName   string    `json:"model_name"`
    ModelVersion string   `json:"model_version"`
    InputData   string    `json:"input_data"`   
    OutputData  string    `json:"output_data"`  
    CreatedAt   time.Time `json:"created_at"`
    HasFeedback bool      `json:"has_feedback"`
}
