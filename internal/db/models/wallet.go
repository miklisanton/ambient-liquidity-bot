package models

import "time"

type Wallet struct {
    ID               int        `db:"id"`
    UserID           int64        `db:"user_id"`
    Address          string     `db:"address"`
    CreatedAt        *time.Time  `db:"created_at"`
}

