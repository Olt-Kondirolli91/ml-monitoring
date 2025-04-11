CREATE TABLE IF NOT EXISTS feedback (
    id UUID PRIMARY KEY,
    inference_id UUID NOT NULL,
    feedback_data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_inference
        FOREIGN KEY (inference_id)
            REFERENCES inferences(id)
);

CREATE INDEX IF NOT EXISTS index_feedback_inference_id
    ON feedback (inference_id);