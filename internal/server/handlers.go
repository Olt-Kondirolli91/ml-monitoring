package server

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/Olt-Kondirolli91/ml-monitoring/internal/models"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}

// handleCreateInference expects a JSON body like:
// {
//   "model_name": "string",
//   "model_version": "string",
//   "input_data": {"some":"input"},
//   "output_data": {"some":"output"}
// }
func (s *Server) handleCreateInference(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ModelName    string      `json:"model_name"`
        ModelVersion string      `json:"model_version"`
        InputData    interface{} `json:"input_data"`
        OutputData   interface{} `json:"output_data"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Generate an ID for the inference
    infID := uuid.New().String()

    // Convert input/output_data to raw JSON string
    inputBytes, _ := json.Marshal(req.InputData)
    outputBytes, _ := json.Marshal(req.OutputData)

    inf := models.Inference{
        ID:           infID,
        ModelName:    req.ModelName,
        ModelVersion: req.ModelVersion,
        InputData:    string(inputBytes),  // store as JSON string
        OutputData:   string(outputBytes), // store as JSON string
        HasFeedback:  false,
    }

    ctx := context.Background()
    if err := s.InferenceRepo.InsertInference(ctx, inf); err != nil {
        log.Printf("Error inserting inference: %v\n", err)
        http.Error(w, "Failed to insert inference", http.StatusInternalServerError)
        return
    }

    // Return the new inference ID
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"inference_id": infID})
}

// handleGetInference retrieves a single inference by ID
func (s *Server) handleGetInference(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    infID := vars["id"]

    ctx := context.Background()
    inf, err := s.InferenceRepo.GetInferenceByID(ctx, infID)
    if err != nil {
        log.Printf("Error getting inference by ID: %v\n", err)
        http.Error(w, "Inference not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(inf)
}

// handleCreateFeedback expects a JSON body like:
// {
//   "feedback_data": {"corrected_output": "foo"}
// }
func (s *Server) handleCreateFeedback(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    infID := vars["id"]

    // parse JSON
    var body struct {
        FeedbackData interface{} `json:"feedback_data"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    fbID := uuid.New().String()
    feedbackBytes, _ := json.Marshal(body.FeedbackData)

    fb := models.Feedback{
        ID:           fbID,
        InferenceID:  infID,
        FeedbackData: string(feedbackBytes),
    }

    ctx := context.Background()
    // Insert feedback
    if err := s.FeedbackRepo.InsertFeedback(ctx, fb); err != nil {
        log.Printf("Error inserting feedback: %v\n", err)
        http.Error(w, "Failed to insert feedback", http.StatusInternalServerError)
        return
    }

    // Mark has_feedback = true
    if err := s.InferenceRepo.UpdateHasFeedback(ctx, infID, true); err != nil {
        log.Printf("Error updating inference has_feedback: %v\n", err)
        http.Error(w, "Failed to update inference feedback status", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"feedback_id": fbID})
}

// handleGetFeedback retrieves all feedback for a given inference
func (s *Server) handleGetFeedback(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    infID := vars["id"]

    ctx := context.Background()
    feedbacks, err := s.FeedbackRepo.GetFeedbackByInferenceID(ctx, infID)
    if err != nil {
        log.Printf("Error getting feedback by inferenceID: %v\n", err)
        http.Error(w, "Feedback not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(feedbacks)
}
