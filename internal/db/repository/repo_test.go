package repository

import (
	"ambient/internal/config"
	"ambient/internal/db/drivers"
	"ambient/internal/db/models"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var walletRepo *WalletRepo
var userRepo *UserRepo
var positionRepo *PositionRepo

func findProjectRoot() (string, error) {
	// Get the current working directory of the test
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse upwards to find the project root (where .env file is located)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); os.IsNotExist(err) {
			parent := filepath.Dir(dir)
			if parent == dir {
				// Reached the root of the filesystem, and .env wasn't found
				return "", os.ErrNotExist
			}
			dir = parent
		} else {
			return dir, nil
		}
	}
}

func TestMain(m *testing.M) {
	root, err := findProjectRoot()
	fmt.Println("Root: ", root)
	if err != nil {
		panic("Error finding project root: " + err.Error())
	}

	err = os.Chdir(root)
	if err != nil {
		panic(err)
	}

    cfg, err := config.NewConfig("config.yaml")
    if err != nil {
        fmt.Println(err)
        return
    }

	connectionURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host+":"+cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSL)
    db, err := drivers.Connect(connectionURL)
    if err != nil {
        fmt.Println(err)
        return
    }
    walletRepo = NewWalletRepo(db)
    userRepo = NewUserRepo(db)
    positionRepo = NewPositionRepo(db)
    m.Run()
}

func TestSaveUser(t *testing.T) {
    user := &models.User{
        ChatID:    13456,
    }
    err := userRepo.Save(context.Background(), user)
    if err != nil {
        t.Errorf("failed to save user: %v", err)
    }
}

func TestGetUserByChatID(t *testing.T) {
    user, err := userRepo.GetByChatID(context.Background(), 13456)
    if err != nil {
        t.Errorf("failed to get user by chat_id: %v", err)
    }
    fmt.Println(user)
}

func TestSaveWallet(t *testing.T) {
    wallet := &models.Wallet{
        UserID: 13456,
        Address: "0x2e1389F741dD0651ea8D235F0460d6C5d15cDFFe",
    }
    err := walletRepo.Save(context.Background(), wallet)
    if err != nil {
        t.Errorf("failed to save wallet: %v", err)
    }
}

func TestSavePosition(t *testing.T) {
    time := time.Now()
    position := &models.Position{
        PositionID: "pos_99318d6f4d380e8e75671919a106cc71e2496d250555862ae8528777a90d63d4",
        WalletID: 1,
        CreatedAt: &time,
        AskTick: 123121,
        BidTick: 123121,
    }
    err := positionRepo.Save(context.Background(), position)
    if err != nil {
        t.Errorf("failed to save position: %v", err)
    }
}

func TestSavePosition2(t *testing.T) {
    time := time.Now().Add(-time.Hour)
    position := &models.Position{
        PositionID: "pos_99318d6f4d380e8e75671919a106cc71e2496d250555862ae8528777a90d63d5",
        WalletID: 1,
        CreatedAt: &time,
        AskTick: 213,
        BidTick: 21323,
    }
    err := positionRepo.Save(context.Background(), position)
    if err != nil {
        t.Errorf("failed to save position: %v", err)
    }
}

func TestSetNotified(t *testing.T) {
    err := positionRepo.SetNotifiedAt(context.Background(), "pos_1")
    if err != nil {
        t.Errorf("failed to set notified_at: %v", err)
    }
}

func TestSetInactive(t *testing.T) {
    err := positionRepo.SetInactive(context.Background(), "pos_1")
    if err != nil {
        t.Errorf("failed to set inactive: %v", err)
    }
}

func TestGetAllPositions(t *testing.T) {
    positions, err := positionRepo.GetAll(context.Background())
    if err != nil {
        t.Errorf("failed to get all positions: %v", err)
    }
    fmt.Println(positions)
}

func TestGetPositionsByUserID(t *testing.T) {
    positions, err := positionRepo.GetByUserId(context.Background(), 13456)
    if err != nil {
        t.Errorf("failed to get positions by user_id: %v", err)
    }
    fmt.Println(positions)
}

func TestGetAllJoined(t *testing.T) {
    positions, err := positionRepo.GetAllJoined(context.Background())
    if err != nil {
        t.Errorf("failed to get all joined: %v", err)
    }
    for _, p := range positions {
        t.Log(p)
    }    
}
