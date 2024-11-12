-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE position(
    position_id char(68) PRIMARY KEY,
    wallet_id INT NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    notified_at TIMESTAMP,
    max_price DOUBLE PRECISION NOT NULL,
    min_price DOUBLE PRECISION NOT NULL,
    CONSTRAINT fk_wallet
        FOREIGN KEY (wallet_id)
        REFERENCES wallet(id)
        ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE position;
-- +goose StatementEnd
