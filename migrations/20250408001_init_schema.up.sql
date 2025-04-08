-- Enable the pgcrypto extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create the 'inferences' table
CREATE TABLE IF NOT EXISTS inferences (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    model_name TEXT NOT NULL,
    model_version TEXT NOT NULL,
    input_data JSONB NOT NULL,
    output_data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    has_feedback BOOLEAN NOT NULL DEFAULT FALSE
);

-- Composite index on (model_name, model_version)
CREATE INDEX IF NOT EXISTS index_inferences_model_name_version
    ON inferences (model_name, model_version);

-- Create the 'feedback' table
CREATE TABLE IF NOT EXISTS feedback (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    inference_id UUID NOT NULL,
    feedback_data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_inference
        FOREIGN KEY(inference_id)
            REFERENCES inferences(id)
);

-- Index on feedback.inference_id
CREATE INDEX IF NOT EXISTS index_feedback_inference_id
    ON feedback (inference_id);
