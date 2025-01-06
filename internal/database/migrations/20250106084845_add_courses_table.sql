-- +goose Up
-- +goose StatementBegin
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Unique identifier using UUID
    title VARCHAR(255) NOT NULL,                   -- Course name
    image_src TEXT,                             -- Course description
    created_at TIMESTAMP DEFAULT NOW(),           -- Creation timestamp
    updated_at TIMESTAMP DEFAULT NOW()            -- Last updated timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS courses;
-- +goose StatementEnd
