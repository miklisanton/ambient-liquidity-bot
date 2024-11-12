package models

type UserWalletPosition struct {
    Address     string  `db:"address"`
    ChatID      int64   `db:"user_id"`
    Position    
}
