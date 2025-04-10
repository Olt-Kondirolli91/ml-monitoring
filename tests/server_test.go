package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/Olt-Kondirolli91/ml-monitoring/internal/server"
    "github.com/gorilla/mux"
    "github.com/google/uuid"
)

func setupMockServer() *server.Server {
    // Create our in-memory mocks
    infRepo := NewMockInferenceRepo()
    fbRepo := NewMockFeedbackRepo()

    s := &server.Server{
        InferenceRepo: infRepo,
        FeedbackRepo:  fbRepo,
        Router:        mux.NewRouter(),
    }
    s.Routes() // call the public method that sets up the routes
    return s
}

func TestCreateInference_Success(t *testing.T) {
    s := setupMockServer()

    body := []byte(`{
        "model_name":"test_model",
        "model_version":"1.0",
        "input_data": {"foo": "bar"},
        "output_data": {"prediction": "xyz"}
    }`)

    req, _ := http.NewRequest("POST", "/inferences", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    s.Router.ServeHTTP(rr, req)
    if rr.Code != http.StatusCreated {
        t.Errorf("Expected 201 Created, got %d", rr.Code)
    }

    var resp map[string]string
    json.NewDecoder(rr.Body).Decode(&resp)
    if resp["inference_id"] == "" {
        t.Error("Expected inference_id in response, got empty")
    }
}

func TestCreateInference_BadJSON(t *testing.T) {
    s := setupMockServer()

    // Missing closing brace
    body := []byte(`{"model_name": "test_model", "model_version":"1.0"`)
    req, _ := http.NewRequest("POST", "/inferences", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    s.Router.ServeHTTP(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Errorf("Expected 400 Bad Request for invalid JSON, got %d", rr.Code)
    }
}

func TestGetInference_NotFound(t *testing.T) {
    s := setupMockServer()

    infID := uuid.New().String() // random non-existent
    req, _ := http.NewRequest("GET", "/inferences/"+infID, nil)
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Errorf("Expected 404 Not Found, got %d", rr.Code)
    }
}

func TestCreateFeedback_Success(t *testing.T) {
    s := setupMockServer()

    // 1) create an inference first, so there's something to attach feedback to
    createBody := []byte(`{
        "model_name":"test_model",
        "model_version":"1.0",
        "input_data": {},
        "output_data": {}
    }`)
    req, _ := http.NewRequest("POST", "/inferences", bytes.NewBuffer(createBody))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    var infResp map[string]string
    json.NewDecoder(rr.Body).Decode(&infResp)
    infID := infResp["inference_id"]

    // 2) post feedback
    fbBody := []byte(`{"feedback_data": {"corrected":"value"}}`)
    fbReq, _ := http.NewRequest("POST", "/inferences/"+infID+"/feedback", bytes.NewBuffer(fbBody))
    fbReq.Header.Set("Content-Type", "application/json")
    fbRR := httptest.NewRecorder()
    s.Router.ServeHTTP(fbRR, fbReq)

    if fbRR.Code != http.StatusCreated {
        t.Errorf("Expected 201 Created, got %d", fbRR.Code)
    }
}

func TestCreateFeedback_NonExistent(t *testing.T) {
    s := setupMockServer()

    fbBody := []byte(`{"feedback_data": {"corrected":"value"}}`)
    req, _ := http.NewRequest("POST", "/inferences/bad-inf-id/feedback", bytes.NewBuffer(fbBody))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    s.Router.ServeHTTP(rr, req)
    // Our mock might handle "bad-inf-id" as not found, so we'd expect 500 or 404 depending on your code logic.
    // In the real code, we do an UpdateHasFeedback => "inference not found".
    // We'll assume you return 500 (since no row was found in DB). If you handle that as 404, adjust here.
    if rr.Code != http.StatusInternalServerError {
        t.Errorf("Expected 500 or 404, got %d", rr.Code)
    }
}

func TestGetFeedback_Empty(t *testing.T) {
    s := setupMockServer()

    // create an inference
    createBody := []byte(`{"model_name":"x","model_version":"y","input_data":{},"output_data":{}}`)
    req, _ := http.NewRequest("POST", "/inferences", bytes.NewBuffer(createBody))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    var infResp map[string]string
    json.NewDecoder(rr.Body).Decode(&infResp)
    infID := infResp["inference_id"]

    // now get feedback => should be empty array
    getReq, _ := http.NewRequest("GET", "/inferences/"+infID+"/feedback", nil)
    getRR := httptest.NewRecorder()
    s.Router.ServeHTTP(getRR, getReq)

    if getRR.Code != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", getRR.Code)
    }

    var data []map[string]interface{}
    json.NewDecoder(getRR.Body).Decode(&data)
    if len(data) != 0 {
        t.Errorf("Expected an empty array, got %d items", len(data))
    }
}
