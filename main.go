package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

const (
	GeckoURL = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&page=1"
)

type CoinInfo struct {
	ID                           string      `json:"id"`
	Symbol                       string      `json:"symbol"`
	Name                         string      `json:"name"`
	Image                        string      `json:"image"`
	CurrentPrice                 float64     `json:"current_price"`
	MarketCap                    int64       `json:"market_cap"`
	MarketCapRank                int         `json:"market_cap_rank"`
	FullyDilutedValuation        float64     `json:"fully_diluted_valuation"`
	TotalVolume                  float64     `json:"total_volume"`
	High24H                      float64     `json:"high_24h"`
	Low24H                       float64     `json:"low_24h"`
	PriceChange24H               float64     `json:"price_change_24h"`
	PriceChangePercentage24H     float64     `json:"price_change_percentage_24h"`
	MarketCapChange24H           float64     `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64     `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64     `json:"circulating_supply"`
	TotalSupply                  float64     `json:"total_supply"`
	MaxSupply                    float64     `json:"max_supply"`
	Ath                          float64     `json:"ath"`
	AthChangePercentage          float64     `json:"ath_change_percentage"`
	AthDate                      time.Time   `json:"ath_date"`
	Atl                          float64     `json:"atl"`
	AtlChangePercentage          float64     `json:"atl_change_percentage"`
	AtlDate                      time.Time   `json:"atl_date"`
	Roi                          interface{} `json:"roi"`
	LastUpdated                  time.Time   `json:"last_updated"`
}

// GetMarketData retrieves the array of current prices from coingecko
func GetMarketData() ([]CoinInfo, error) {
	var prices []CoinInfo

	req, err := http.NewRequest("GET", GeckoURL, nil)
	if err != nil {
		return prices, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("accept", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return prices, err
	}

	if resp.StatusCode == 429 {
		fmt.Println("Being rate limited by coingecko")
		return prices, nil
	}

	results, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return prices, err
	}

	err = json.Unmarshal(results, &prices)
	if err != nil {
		fmt.Printf(resp.Status)
		return prices, err
	}

	return prices, nil
}

// store puts the coins values into redis
func store(client *redis.Client, coin CoinInfo) {
	client.Set(ctx, coin.ID+"#Symbol", coin.Symbol, time.Minute)
	client.Set(ctx, coin.ID+"#Name", coin.Name, time.Minute)
	client.Set(ctx, coin.ID+"#Image", coin.Image, time.Minute)
	client.Set(ctx, coin.ID+"#CurrentPrice", coin.CurrentPrice, time.Minute)
	client.Set(ctx, coin.ID+"#MarketCap", coin.MarketCap, time.Minute)
	client.Set(ctx, coin.ID+"#MarketCapRank", coin.MarketCapRank, time.Minute)
	client.Set(ctx, coin.ID+"#FullyDilutedValuation", coin.FullyDilutedValuation, time.Minute)
	client.Set(ctx, coin.ID+"#TotalVolume", coin.TotalVolume, time.Minute)
	client.Set(ctx, coin.ID+"#High24H", coin.High24H, time.Minute)
	client.Set(ctx, coin.ID+"#Low24H", coin.Low24H, time.Minute)
	client.Set(ctx, coin.ID+"#PriceChange24H", coin.PriceChange24H, time.Minute)
	client.Set(ctx, coin.ID+"#PriceChangePercentage24H", coin.PriceChangePercentage24H, time.Minute)
	client.Set(ctx, coin.ID+"#MarketCapChange24H", coin.MarketCapChange24H, time.Minute)
	client.Set(ctx, coin.ID+"#MarketCapChangePercentage24H", coin.MarketCapChangePercentage24H, time.Minute)
	client.Set(ctx, coin.ID+"#CirculatingSupply", coin.CirculatingSupply, time.Minute)
	client.Set(ctx, coin.ID+"#TotalSupply", coin.TotalSupply, time.Minute)
	client.Set(ctx, coin.ID+"#MaxSupply", coin.MaxSupply, time.Minute)
	client.Set(ctx, coin.ID+"#Ath", coin.Ath, time.Minute)
	client.Set(ctx, coin.ID+"#AthChangePercentage", coin.AthChangePercentage, time.Minute)
	client.Set(ctx, coin.ID+"#AthDate", coin.AthDate, time.Minute)
	client.Set(ctx, coin.ID+"#Atl", coin.Atl, time.Minute)
	client.Set(ctx, coin.ID+"#AtlChangePercentage", coin.AtlChangePercentage, time.Minute)
	client.Set(ctx, coin.ID+"#AtlDate", coin.AtlDate, time.Minute)
	client.Set(ctx, coin.ID+"#Roi", coin.Roi, time.Minute)
	client.Set(ctx, coin.ID+"#LastUpdated", coin.LastUpdated, time.Minute)
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	for {

		coinsData, err := GetMarketData()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, coin := range coinsData {
			go store(rdb, coin)
		}

		time.Sleep(1 * time.Second)
	}
}
