package models

type PositionWallet struct {
    Position
    Address          string `db:"address"`
}
