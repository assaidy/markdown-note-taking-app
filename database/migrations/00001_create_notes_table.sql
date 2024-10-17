-- +goose Up
-- +goose StatementBegin
CREATE TABLE notes (
    id         INTEGER      PRIMARY KEY AUTOINCREMENT,
    title      VARCHAR(225) NOT NULL,
    content    TEXT         NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notes;
-- +goose StatementEnd
