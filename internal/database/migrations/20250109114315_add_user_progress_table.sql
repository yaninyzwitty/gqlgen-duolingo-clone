-- +goose Up
-- +goose StatementBegin
-- Creating the "user_progress" table with a foreign key to "course"
CREATE TABLE IF NOT EXISTS "user_progress" (
    user_id UUID PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    active_course_id UUID,
    hearts INT NOT NULL DEFAULT 0,
    points INT NOT NULL DEFAULT 0,
    FOREIGN KEY (active_course_id) REFERENCES courses(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Create an index for faster querying by active_course_id
CREATE INDEX IF NOT EXISTS idx_active_course_id ON "user_progress"(active_course_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Dropping the index and table if they exist
DROP INDEX IF EXISTS idx_active_course_id;
DROP TABLE IF EXISTS "user_progress";
-- +goose StatementEnd