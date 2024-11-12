package services

import (
	"ambient/internal/db/models"
	"ambient/internal/db/repository"
	"context"
	"fmt"
)

type WalletService struct {
    walletRepo *repository.WalletRepo
}

func NewWalletService(walletRepo *repository.WalletRepo) *WalletService {
    return &WalletService{
        walletRepo: walletRepo,
    }
}

func (w *WalletService) Save(ctx context.Context, wallet *models.Wallet) error {
    if wallet.UserID == 0 {
        return fmt.Errorf("user_id is required")
    } 
    if wallet.Address == "" {
        return fmt.Errorf("address is required")
    }
    return w.walletRepo.Save(ctx, wallet)
}

func (w *WalletService) GetByUserID(ctx context.Context, userID int64) ([]models.Wallet, error) {
    return w.walletRepo.GetByUserID(ctx, userID)
}

func (w *WalletService) GetAll(ctx context.Context) ([]models.Wallet, error) {
    return w.walletRepo.GetAll(ctx)
}

func (w *WalletService) GetByAddress(ctx context.Context, address string) (*models.Wallet, error) {
    return w.walletRepo.GetByAddress(ctx, address)
}
