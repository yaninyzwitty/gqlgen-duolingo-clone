-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- Ensure pgcrypto extension is enabled

CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- Generate a UUID by default
    title VARCHAR(255) NOT NULL,                    -- Title of the unit
    description TEXT,                               -- Optional description of the unit
    course_id UUID NOT NULL,                        -- Foreign key to link the unit to a course
    unit_order SERIAL,                              -- Auto-incrementing order value
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- Timestamp when the unit was created
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- Timestamp for last update
    CONSTRAINT fk_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE INDEX idx_course_id ON units(course_id);

-- Create a trigger to automatically update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_units_updated_at
BEFORE UPDATE ON units
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_units_updated_at ON units;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP TABLE IF EXISTS units;
-- +goose StatementEnd