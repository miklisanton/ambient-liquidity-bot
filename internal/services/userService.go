package services

import (
    "ambient/internal/db/models"
    "ambient/internal/db/repository"
    "context"
)

type UserService struct {
    userRepo *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
    return &UserService{
        userRepo: userRepo,
    }
}

func (u *UserService) Save(ctx context.Context, user *models.User) error {
    return u.userRepo.Save(ctx, user)
}

func (u *UserService) GetByChatID(ctx context.Context, chatID int64) (*models.User, error) {
    return u.userRepo.GetByChatID(ctx, chatID)
}
