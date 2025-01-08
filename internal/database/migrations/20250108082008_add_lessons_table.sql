-- +goose Up
-- +goose StatementBegin
CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    unit_id UUID NOT NULL,
    unit_order INT NOT NULL,
    FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lessons;
-- +goose StatementEnd
