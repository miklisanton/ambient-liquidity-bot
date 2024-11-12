package repository

import (
    "fmt"
    "ambient/internal/db/models"
    "context"
    "github.com/jmoiron/sqlx"
)

type WalletRepo struct {
    db *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) *WalletRepo {
    return &WalletRepo{
        db: db,
    }
}

func (w *WalletRepo) Save(ctx context.Context, wallet *models.Wallet) error {
    query := `
    INSERT INTO wallet (user_id, address, created_at)
    VALUES ($1, $2, COALESCE($3, NOW()))`
    _, err := w.db.ExecContext(ctx, query, wallet.UserID, wallet.Address, wallet.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to insert wallet: %w", err)
    }
    return nil
}

func (w *WalletRepo) GetByUserID(ctx context.Context, userID int64) ([]models.Wallet, error) {
    query := `SELECT * FROM wallet WHERE user_id = $1`
    var wallets []models.Wallet
    err := w.db.SelectContext(ctx, &wallets, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get wallets by user_id: %w", err)
    }
    return wallets, nil
}

func (w *WalletRepo) GetAll(ctx context.Context) ([]models.Wallet, error) {
    query := `SELECT * FROM wallet`
    var wallets []models.Wallet
    err := w.db.SelectContext(ctx, &wallets, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get all wallets: %w", err)
    }
    return wallets, nil
}

func (w *WalletRepo) GetByAddress(ctx context.Context, address string) (*models.Wallet, error) {
    query := `SELECT * FROM wallet WHERE address = $1`
    var wallet models.Wallet
    err := w.db.GetContext(ctx, &wallet, query, address)
    if err != nil {
        return nil, fmt.Errorf("failed to get wallet by address: %w", err)
    }
    return &wallet, nil
}

