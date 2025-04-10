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

-- Seed some test data
-- -------------------------------------------------------
INSERT INTO inferences (id, model_name, model_version, input_data, output_data, has_feedback)
VALUES
('c3c4c350-7d8a-4b02-816c-245ced77ff01', 'seed_model', '1.0.0', '{"example":"input"}', '{"example":"output"}', false),
('11111111-2222-3333-4444-555555555555', 'seed_model', '1.0.0', '{"key":"val"}', '{"prediction":"some output"}', true);

-- One row in feedback referencing the second inference
INSERT INTO feedback (id, inference_id, feedback_data)
VALUES (
    '66666666-7777-8888-9999-000000000000',
    '11111111-2222-3333-4444-555555555555',
    '{"corrected_output":"updated result"}'
);