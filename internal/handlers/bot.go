package handlers

import (
	"ambient/internal/db/models"
	"ambient/internal/services"
	"ambient/internal/utils"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
    api *tgbotapi.BotAPI
    userS *services.UserService
    walletS *services.WalletService
    positionS *services.PositionService
    ambientS *services.AmbientService
    binanceS *services.BinanceService
    newWalletCh chan bool
    notificationCh chan utils.INotification
}



func NewBot(
    api *tgbotapi.BotAPI,
    userS *services.UserService,
    walletS *services.WalletService,
    positionS *services.PositionService,
    ambientS *services.AmbientService,
    binanceS *services.BinanceService) *Bot {
    return &Bot{
        api: api,
        userS: userS,
        walletS: walletS,
        positionS: positionS,
        ambientS: ambientS,
        binanceS: binanceS,
        newWalletCh: make(chan bool),
        notificationCh: make(chan utils.INotification, 512),
    }
}

func (b *Bot) Start() {
    log.Println("Starting bot")
    go func() {
        // Set available commands
        commands := []tgbotapi.BotCommand{
            {Command: "start", Description: "Start interacting with the bot"},
            {Command: "a", Description: "/a <wallet_address> - Add a wallet"},
            {Command: "l", Description: "/l - List all wallets"},
            {Command: "p", Description: "/p <wallet_address|null> - List all positions for a wallet"},
        }

        setCommandsConfig := tgbotapi.NewSetMyCommands(commands...)

        _, err := b.api.Request(setCommandsConfig)
        if err != nil {
            log.Panic(err)
        }

        log.Println("Commands set successfully!")

        updates := tgbotapi.NewUpdate(0)
        updates.Timeout = 60
        updatesChan := b.api.GetUpdatesChan(updates)
        if err != nil {
            log.Fatalf("Error getting updates: %s", err)
        }
        for update := range updatesChan {
            if update.Message == nil {
                continue
            }
            chatID := update.Message.Chat.ID
            switch update.Message.Command() {
            case "start":
                if err := b.HandleStart(update.Message); err != nil {
                    log.Printf("Error handling start: %s", err)
                }
            case "a":
                if err := b.HandleAddWallet(update.Message); err != nil {
                    log.Printf("Error handling add wallet: %s", err)
                }
            case "l":
                if err := b.HandleListWallets(update.Message); err != nil {
                    log.Printf("Error handling list wallets: %s", err)
                }
            case "p":
                if err := b.HandleListPositions(update.Message); err != nil {
                    log.Printf("Error handling list positions: %s", err)
                }
            default:
                b.SendMessage(chatID, "Invalid command")
            }
        }
    }()
    b.Monitor()
    b.TrackPrice()
}

func (b *Bot)HandleStart(msg *tgbotapi.Message) error {
    log.Println("Handling start")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    user := &models.User{
        ChatID: msg.Chat.ID,
    }
    b.userS.Save(ctx, user)
    b.SendMessage(msg.Chat.ID, "Welcome to Ambient Bot!")
    return nil
}

//
// /a
func (b *Bot)HandleAddWallet(msg *tgbotapi.Message) error {
    log.Println("Handling add wallet")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    wallet := &models.Wallet{
        UserID: msg.Chat.ID,
        Address: msg.CommandArguments(),
    }
    b.walletS.Save(ctx, wallet)
    return nil
}

//
// /l
func (b *Bot)HandleListWallets(msg *tgbotapi.Message) error {
    log.Println("Handling list wallets")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    wallets, err := b.walletS.GetByUserID(ctx, msg.Chat.ID)
    if err != nil {
        return fmt.Errorf("Error getting wallets: %s", err)
    }
    if len(wallets) == 0 {
        b.SendMessage(msg.Chat.ID, "No wallets found")
    } else {
        for _, wallet := range wallets {
            b.SendMessage(msg.Chat.ID, wallet.Address)
        }
    }
    return nil
}

// 
// /p <wallet_address|null>
func (b *Bot) HandleListPositions(msg *tgbotapi.Message) error {
    log.Println("Handling list positions")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    walletAddress := msg.CommandArguments()
    if walletAddress == "" {
        wallets, err := b.walletS.GetByUserID(ctx, msg.Chat.ID)
        if err != nil {
            return fmt.Errorf("Error getting wallets: %s", err)
        }
        if len(wallets) == 0 {
            b.SendMessage(msg.Chat.ID, "No wallets found")
        } else {
            for _, wallet := range wallets {
                positions, err := b.positionS.GetByWalletID(ctx, wallet.ID)
                if err != nil {
                    return fmt.Errorf("Error getting positions: %s", err)
                }
                if len(positions) == 0 {
                    b.SendMessage(msg.Chat.ID, "wallet: " + wallet.Address + "\nNo positions found")
                } else {
                    for _, position := range positions {
                        b.SendMessage(msg.Chat.ID, fmt.Sprintf("wallet: %s\nposition: %s\nmax_price: %f\nmin_price: %f\nactive: %t", wallet.Address, position.PositionID, position.MaxPrice, position.MinPrice, position.Active))
                    }
                }
            }
        }
    } else {
        wallet, err := b.walletS.GetByAddress(ctx, walletAddress)
        if err != nil {
            return fmt.Errorf("Error getting wallet: %s", err)
        }
        if wallet == nil {
            b.SendMessage(msg.Chat.ID, "Wallet not found")
        } else {
            positions, err := b.positionS.GetByWalletID(ctx, wallet.ID)
            if err != nil {
                return fmt.Errorf("Error getting positions: %s", err)
            }
            if len(positions) == 0 {
                b.SendMessage(msg.Chat.ID, "No positions found")
            } else {
                for _, position := range positions {
                    b.SendMessage(msg.Chat.ID, fmt.Sprintf("wallet: %s\nposition: %s\nmax_price: %f\nmin_price: %f\nactive: %t", wallet.Address, position.PositionID, position.MaxPrice, position.MinPrice, position.Active))
                }
            }
        }
    }
    return nil
}

func (b *Bot) NotifyUser(msg *tgbotapi.Message) error {
    log.Println("Notifying user")
    return nil
}

func (bot *Bot) SendMessage(chatID int64, text string) int {
	msg := tgbotapi.NewMessage(chatID, text)
	msgSent, err := bot.api.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return 0
	}

	return msgSent.MessageID
}

func (b *Bot) Monitor() { 
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        for {
            select {
            case <-ticker.C:
                if err := b.CheckWallets(); err != nil {
                    log.Printf("Error checking wallets: %s", err)
                }
            case <-b.newWalletCh:
                if err := b.CheckWallets(); err != nil {
                    log.Printf("Error checking wallets: %s", err)
                }
            }
        }
    } ()
}

func (b *Bot) CheckWallets() error {
    log.Println("Checking all wallets for new positions")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    wallets, err := b.walletS.GetAll(ctx)
    if err != nil {
        return fmt.Errorf("Error getting wallets: %s", err)
    }
    log.Printf("Found %d wallets", len(wallets))
    wg := sync.WaitGroup{}
    for _, wallet := range wallets {
        wg.Add(1)
        go func(wallet models.Wallet) {
            defer wg.Done()
            if err := b.UpdatePositions(ctx, &wallet); err != nil {
                log.Printf("%s: Error updating positions: %s", wallet.Address, err)
            }
        }(wallet)
    }
    wg.Wait()
    log.Println("All wallets checked")
    return nil
}

func (b *Bot) UpdatePositions(ctx context.Context, wallet *models.Wallet) error {
    select {
    case <-ctx.Done():
        log.Printf("Context done")
        return nil
    default:
        log.Printf("Checking positions for wallet %s", wallet.Address)
        activePositions, err := b.ambientS.GetUserPools(wallet.Address)
        if err != nil {
            return fmt.Errorf("error getting user pools: %s", err)
        }
        savedPositions, err := b.positionS.GetByWalletID(ctx, wallet.ID)
        if err != nil {
            return fmt.Errorf("error getting saved positions: %s", err)
        }
        savedMap := make(map[string]models.UserWalletPosition)
        for _, savedPos := range savedPositions {
            // Save position to map
            savedMap[savedPos.PositionID] = savedPos
            // Set flag if position is inactive
            if _, ok := activePositions[savedPos.PositionID]; ok && savedPos.Active {
                continue
            } else {
                err := b.positionS.SetInactive(ctx, savedPos.PositionID)
                b.notificationCh <- &utils.LiquidationNotification{
                    Position: savedPos,
                }

                if err != nil {
                    return fmt.Errorf("Error setting position inactive: %s", err)
                }
            }
        }
        // Add new positions
        for _, activePos := range activePositions {
            if _, ok := savedMap[activePos.PositionId]; ok {
                continue
            }
            t := time.Unix(activePos.LatestUpdateTime, 0)
            err := b.positionS.Save(ctx, &models.Position{
                PositionID: activePos.PositionId,
                WalletID: wallet.ID,
                CreatedAt: &t,
                MaxPrice: utils.TickToPrice(activePos.BidTick, 18, 6),
                MinPrice: utils.TickToPrice(activePos.AskTick, 18, 6),
            })
            if err != nil {
                return fmt.Errorf("Error saving position: %s", err)
            }
        }
    }
    return nil
}

func (b *Bot) TrackPrice() {
    // Start binance service
    b.binanceS.Connect()
    b.binanceS.Listen()

    ticker := time.NewTicker(15 * time.Second)

    mutex := sync.Mutex{}

    var positions []models.UserWalletPosition
    var err error

    go func() {
        // Get all active positions when starting
        log.Println("Getting all active positions")
        mutex.Lock()
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        positions, err = b.positionS.GetAllJoined(ctx)
        if err != nil {
            log.Printf("Error getting all positions: %s", err)
        }
        log.Printf("Found %d active positions", len(positions))
        mutex.Unlock()
        // Get all active positions every tick
        for {
            select {
            case <-ticker.C:
                // Get all active positions
                log.Println("Getting all active positions")
                mutex.Lock()
                
                ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
                positions, err = b.positionS.GetAllJoined(ctx)
                cancel()
                if err != nil {
                    log.Printf("Error getting all positions: %s", err)
                    continue
                }
                mutex.Unlock()
            }
        }
    } ()

    go func() {
        for price := range b.binanceS.Out {
            mutex.Lock()
            for _, position := range positions {
                if position.MaxPrice < price * 1.01 || position.MinPrice > price * 0.99 {
                    b.notificationCh <- &utils.Notification{
                        Position: position,
                        Price: price,
                    }
                    if err := b.positionS.SetNotifiedAt(context.Background(), position.PositionID); err != nil {
                        log.Printf("Error setting notified_at: %s", err)
                    }
                }
            }
            mutex.Unlock()
        }
    } ()

    for notification := range b.notificationCh {
        log.Printf(notification.GetMessage())
        b.SendMessage(notification.GetChatID(), notification.GetMessage())
    }
}
