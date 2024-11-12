-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE wallet(
    id SERIAL PRIMARY KEY,
    address CHAR(42) NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(chat_id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE wallet;
-- +goose StatementEnd
