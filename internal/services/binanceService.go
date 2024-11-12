package services

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

type BinanceService struct {
    BaseUrl string
    Conn *websocket.Conn
    Out chan float64
}

func NewBinanceService() *BinanceService {
    return &BinanceService{
        BaseUrl: "wss://stream.binance.com:443/ws/ethusdt@ticker",
        Out: make(chan float64),
    }
}

func (b *BinanceService) Connect() error {
    conn, _, err := websocket.DefaultDialer.Dial(b.BaseUrl, nil)
    if err != nil {
        return err
    }
    b.Conn = conn
    return nil
}

func (b *BinanceService) Listen() {
    go func() {
        for {
            _, message, err := b.Conn.ReadMessage()
            if err != nil {
                return
            }
            var data struct {
                Price string `json:"b"`
                B string `json:"B"`
            }
            if err := json.Unmarshal(message, &data); err != nil {
                log.Println(err)
                return
            }
            price, err := strconv.ParseFloat(data.Price, 32)
            if err != nil {
                log.Println(err)
                return
            }
            b.Out <- price
        }
    }()
}

