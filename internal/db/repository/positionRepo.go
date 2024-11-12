package repository

import (
    "fmt"
    "ambient/internal/db/models"
    "context"
    "github.com/jmoiron/sqlx"
)

type PositionRepo struct {
    db *sqlx.DB
}

func NewPositionRepo(db *sqlx.DB) *PositionRepo {
    return &PositionRepo{
        db: db,
    }
}

func (p *PositionRepo) Save(ctx context.Context, position *models.Position) error {
    query := `
    INSERT INTO position (position_id, wallet_id, created_at, max_price, min_price)
    VALUES ($1, $2, $3, $4, $5)`
    _, err := p.db.ExecContext(
        ctx,
        query,
        position.PositionID,
        position.WalletID,
        position.CreatedAt,
        position.MaxPrice,
        position.MinPrice)
    if err != nil {
        return fmt.Errorf("failed to insert position: %w", err)
    }
    return nil
}

func (p *PositionRepo) GetAll(ctx context.Context) ([]models.Position, error) {
    query := `SELECT * FROM position WHERE active = true`
    var positions []models.Position
    err := p.db.SelectContext(ctx, &positions, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get all positions: %w", err)
    }
    return positions, nil
}

func (p *PositionRepo) GetAllJoined(ctx context.Context) ([]models.UserWalletPosition, error) {
    query := `
    SELECT p.*, w.address, w.user_id
    FROM position p
    JOIN wallet w ON p.wallet_id = w.id
    WHERE p.active = true
    AND (NOW() - p.notified_at > INTERVAL '30 minutes') OR p.notified_at IS NULL`

    var positions []models.UserWalletPosition
    err := p.db.SelectContext(ctx, &positions, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get all positions with addresses and chat_ids: %w", err)
    }
    return positions, nil
}

func (p *PositionRepo) GetByWalletID(ctx context.Context, walletID int) ([]models.UserWalletPosition, error) {
    query := `
    SELECT p.*, w.address, w.user_id 
    FROM position p
    JOIN wallet w on p.wallet_id = w.id
    WHERE wallet_id = $1 AND active = true`
    var positions []models.UserWalletPosition
    err := p.db.SelectContext(ctx, &positions, query, walletID)
    if err != nil {
        return nil, fmt.Errorf("failed to get position by wallet_id: %w", err)
    }
    return positions, nil
}

func (p *PositionRepo) GetByUserId(ctx context.Context, userId int) ([]models.PositionWallet, error) {
    query := `
    SELECT p.*, w.address
    FROM position p
    JOIN wallet w ON p.wallet_id = w.id
    WHERE w.user_id = $1 AND p.active = true`

    var positions []models.PositionWallet
    err := p.db.SelectContext(ctx, &positions, query, userId)
    if err != nil {
        return nil, fmt.Errorf("failed to get position by user_id: %w", err)
    }
    return positions, nil
}

func (p *PositionRepo) SetInactive(ctx context.Context, positionID string) error {
    query := `UPDATE position SET active = false WHERE position_id = $1`
    _, err := p.db.ExecContext(ctx, query, positionID)
    if err != nil {
        return fmt.Errorf("failed to set position inactive: %w", err)
    }
    return nil
}

func (p *PositionRepo) SetNotifiedAt(ctx context.Context, positionID string) error {
    query := `UPDATE position SET notified_at = now() WHERE position_id = $1`
    _, err := p.db.ExecContext(ctx, query, positionID)
    if err != nil {
        return fmt.Errorf("failed to set position notified_at: %w", err)
    }
    return nil
}

