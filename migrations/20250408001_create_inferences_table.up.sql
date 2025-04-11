CREATE TABLE IF NOT EXISTS inferences (
    id UUID PRIMARY KEY,
    model_name TEXT NOT NULL,
    model_version TEXT NOT NULL,
    input_data JSONB NOT NULL,
    output_data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    has_feedback BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS index_inferences_model_name_version
    ON inferences (model_name, model_version);