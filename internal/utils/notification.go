package utils

import (
	"ambient/internal/db/models"
	"fmt"
)

type INotification interface {
    GetChatID() int64
    GetMessage() string
}

type Notification struct {
    Position models.UserWalletPosition
    Price float64
}

type LiquidationNotification struct {
    Position models.UserWalletPosition
}

func (n LiquidationNotification) GetChatID() int64 {
    return n.Position.ChatID
}

func (n LiquidationNotification) GetMessage() string {
    return fmt.Sprintf("wallet: %s\npostion id: %s\n has been liquidated.\nmin price: %.2f\nmax price: %.2f", 
        n.Position.Address,
        n.Position.PositionID,
        n.Position.MinPrice,
        n.Position.MaxPrice)
}

func (n Notification) GetChatID() int64 {
    return n.Position.ChatID
}

func (n Notification) GetMessage() string {
    return fmt.Sprintf("Wallet: %s\nPostion id: %s\n is close to liquidation.\nmin price: %.2f\n\nmax price: %.2f\ncurrent price: %.2f", 
        n.Position.Address,
        n.Position.PositionID,
        n.Position.MinPrice,
        n.Position.MaxPrice,
        n.Price)
    }

