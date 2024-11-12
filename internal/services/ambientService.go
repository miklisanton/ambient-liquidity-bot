package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AmbientService struct {
    BaseUrl string
    httpClient *http.Client
}

type Position struct {
    ChainId           string `json:"chainId"`
    Base              string `json:"base"`
    Quote             string `json:"quote"`
    PoolIdx           int    `json:"poolIdx"`
    BidTick           int    `json:"bidTick"`
    AskTick           int    `json:"askTick"`
    IsBid             bool   `json:"isBid"`
    User              string `json:"user"`
    TimeFirstMint     int64  `json:"timeFirstMint"`
    LatestUpdateTime  int64  `json:"latestUpdateTime"`
    LastMintTx        string `json:"lastMintTx"`
    FirstMintTx       string `json:"firstMintTx"`
    PositionType      string `json:"positionType"`
    AmbientLiq        int64  `json:"ambientLiq"`
    ConcLiq           int64  `json:"concLiq"`
    RewardLiq         int64  `json:"rewardLiq"`
    LiqRefreshTime    int64  `json:"liqRefreshTime"`
    AprDuration       int64  `json:"aprDuration"`
    AprPostLiq        float64 `json:"aprPostLiq"`
    AprContributedLiq int64  `json:"aprContributedLiq"`
    AprEst            float64 `json:"aprEst"`
    PositionId        string `json:"positionId"`
}

type Response struct {
    Data []Position `json:"data"`
}

func NewAmbientService() *AmbientService {
    return &AmbientService{
        BaseUrl: "https://ambindexer.net/scroll-gcgo/",
        httpClient: &http.Client{},
    }
}

func (a *AmbientService) GetUserPools(address string) (map[string]Position, error) {
    url := fmt.Sprintf(
        "%suser_positions?chainId=%s&user=%s",
        a.BaseUrl,
        "0x82750",
        address,
    )
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    res, err := a.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    var response Response
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, err
    }
    
    activePositions := make(map[string]Position)
    for _, position := range response.Data {
        if position.ConcLiq > 0 {
            activePositions[position.PositionId] = position
        }
    }
    return activePositions, nil
}

func (a *AmbientService) GetStats(pos Position) (any, error) {
    url := fmt.Sprintf(
        "%suser_positions?chainId=%s&user=%s&poolIdx=%d&quote=%s&base=%s&askTick=%d&bidTick=%d",
        a.BaseUrl,
        "0x82750",
        pos.User,
        pos.PoolIdx,
        pos.Quote,
        pos.Base,
        pos.AskTick,
        pos.BidTick,

    )
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    res, err := a.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    fmt.Println(string(body))
//    var response Response
//    if err := json.Unmarshal(body, &response); err != nil {
//        return nil, err
//    }
//    
//    var activePositions []Position
//    for _, position := range response.Data {
//        if position.ConcLiq > 0 {
//            activePositions = append(activePositions, position)
//        }
//    }
    return nil, nil
}

//func (a *AmbientService) MonitorWallets(newWalletCh chan bool) { 
//    ticker := time.NewTicker(30 * time.Second)
//    for {
//        select {
//        case <-ticker.C:
//        case <-newWalletCh:
//            log.Println("Checking all wallets for new positions")
//            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//            defer cancel()
//
//            wallets, err := .GetAll(ctx)
//            if err != nil {
//                log.Printf("Error getting wallets: %s", err)
//                continue
//            }
//            for _, wallet := range wallets {
//                positions, err := a.GetUserPools(wallet.Address)
//                if err != nil {
//                    log.Printf("Error getting user pools: %s", err)
//                    continue
//                }
//                for _, position := range positions {
//                    err := positionRepo.Save(ctx, &models.Position{
//                        WalletID: wallet.ID,
//                        CreatedAt: time.Unix(position.LatestUpdateTime, 0),
//                        AskTick: position.AskTick,
//                        BidTick: position.BidTick,
//                    })
//                    if err != nil {
//                        log.Printf("Error saving position: %s", err)
//                    }
//                }
//            }
//        }
//    }
//}
