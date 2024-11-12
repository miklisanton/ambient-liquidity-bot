package services

import (
	"ambient/internal/db/models"
	"ambient/internal/db/repository"
	"context"
)

type PositionService struct {
    positionRepo *repository.PositionRepo
}

func NewPositionService(positionRepo *repository.PositionRepo) *PositionService {
    return &PositionService{
        positionRepo: positionRepo,
    }
}

func (p *PositionService) Save(ctx context.Context, position *models.Position) error {
    return p.positionRepo.Save(ctx, position)
}

func (p *PositionService) GetAll(ctx context.Context) ([]models.Position, error) {
    return p.positionRepo.GetAll(ctx)
}

func (p *PositionService) GetByWalletID(ctx context.Context, walletID int) ([]models.UserWalletPosition, error) {
    return p.positionRepo.GetByWalletID(ctx, walletID)
}

func (p *PositionService) GetByUserId(ctx context.Context, userId int) ([]models.PositionWallet, error) {
    return p.positionRepo.GetByUserId(ctx, userId)
}

func (p *PositionService) SetInactive(ctx context.Context, positionID string) error {
    return p.positionRepo.SetInactive(ctx, positionID)
}

func (p *PositionService) SetNotifiedAt(ctx context.Context, positionID string) error {
    return p.positionRepo.SetNotifiedAt(ctx, positionID)
}

func (p *PositionService) GetAllJoined(ctx context.Context) ([]models.UserWalletPosition, error) {
    return p.positionRepo.GetAllJoined(ctx)
}


