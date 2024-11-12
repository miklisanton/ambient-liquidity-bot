package models

import "time"

type Position struct {
    PositionID       string `db:"position_id"`
    WalletID         int    `db:"wallet_id"`
    Active           bool   `db:"active"`
    CreatedAt        *time.Time `db:"created_at"`
    NotifiedAt       *time.Time `db:"notified_at"`
    MaxPrice          float64    `db:"max_price"`
    MinPrice          float64    `db:"min_price"`
}

