package main

import (
	"ambient/internal/config"
	"ambient/internal/db/drivers"
	"ambient/internal/db/repository"
	"ambient/internal/handlers"
	"ambient/internal/services"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)



func main() {
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

    walletRepo := repository.NewWalletRepo(db)
    userRepo := repository.NewUserRepo(db)
    positionRepo := repository.NewPositionRepo(db)

    walletService := services.NewWalletService(walletRepo)
    positionService := services.NewPositionService(positionRepo)
    userService := services.NewUserService(userRepo)
    
    ambientService := services.NewAmbientService()
    binanceService := services.NewBinanceService()
    
    api, err := tgbotapi.NewBotAPI(cfg.Telegram.APIkey)
    bot := handlers.NewBot(api, userService, walletService, positionService, ambientService, binanceService)
    bot.Start()
}
