-- +goose Up
-- +goose StatementBegin
CREATE TYPE challenge_type AS ENUM ('SELECT', 'ASSIST');

CREATE TABLE challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL,
    type challenge_type NOT NULL,
    question TEXT NOT NULL,
    unit_order INT NOT NULL,
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE challenges;
DROP TYPE challenge_type;
-- +goose StatementEnd
