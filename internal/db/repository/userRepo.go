package repository

import (
    "fmt"
    "ambient/internal/db/models"
    "context"
    "github.com/jmoiron/sqlx"
)

type UserRepo struct {
    db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
    return &UserRepo{
        db: db,
    }
}

func (u *UserRepo) Save(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (chat_id, created_at)
        VALUES ($1, COALESCE($2, NOW()))`
    _, err := u.db.ExecContext(ctx, query, user.ChatID, user.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    return nil
}

func (u *UserRepo) GetByChatID(ctx context.Context, chatID int64) (*models.User, error) {
    query := `SELECT * FROM users WHERE chat_id = $1`
    var user models.User
    err := u.db.GetContext(ctx, &user, query, chatID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by chat_id: %w", err)
    }
    return &user, nil
}
