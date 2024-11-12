package models

import (
    "time"
)

type User struct {
    ChatID int64 `db:"chat_id"`
    CreatedAt *time.Time `db:"created_at"`
}

