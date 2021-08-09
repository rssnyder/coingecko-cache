package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx       = context.Background()
	frequency *int
	pages     *int
	expiry    *int
	hostname  *string
	password  *string
	db        *int
	wg        sync.WaitGroup
)

const (
	GeckoURL = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&page=%d"
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

func init() {
	frequency = flag.Int("frequency", 1, "seconds between updates")
	pages = flag.Int("pages", 1, "number of pages (100 coin each) to pull from")
	expiry = flag.Int("expiry", 60, "number of seconds to keep entries in the cache")
	hostname = flag.String("hostname", "localhost:6379", "connection address for redis")
	password = flag.String("password", "", "redis password")
	db = flag.Int("db", 0, "redis db to use")
	flag.Parse()
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     *hostname,
		Password: *password,
		DB:       *db,
	})

	pager := 1

	for {

		coinsData, err := GetMarketData(pager)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(*frequency) * time.Second)
			continue
		}

		for _, coin := range coinsData {
			wg.Add(1)
			go Store(&wg, rdb, coin, time.Duration(*expiry)*time.Second)
		}
		fmt.Println("waiting for storage")
		wg.Wait()

		pager++
		if pager > *pages {
			pager = 1
			time.Sleep(time.Duration(*frequency) * time.Second)
		}
	}
}

// GetMarketData retrieves the array of current prices from coingecko
func GetMarketData(page int) ([]CoinInfo, error) {
	var prices []CoinInfo

	req, err := http.NewRequest("GET", fmt.Sprintf(GeckoURL, page), nil)
	if err != nil {
		return prices, err
	}
	fmt.Printf("Retrived page %d\n", page)

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

// Store puts the coins values into redis
func Store(wg *sync.WaitGroup, client *redis.Client, coin CoinInfo, expiry time.Duration) {
	defer wg.Done()

	err := client.Set(ctx, coin.ID+"#Symbol", coin.Symbol, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#Name", coin.Name, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#Image", coin.Image, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#CurrentPrice", coin.CurrentPrice, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#MarketCap", coin.MarketCap, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#MarketCapRank", coin.MarketCapRank, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#FullyDilutedValuation", coin.FullyDilutedValuation, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#TotalVolume", coin.TotalVolume, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#High24H", coin.High24H, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#Low24H", coin.Low24H, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#PriceChange24H", coin.PriceChange24H, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#PriceChangePercentage24H", coin.PriceChangePercentage24H, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#MarketCapChange24H", coin.MarketCapChange24H, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#MarketCapChangePercentage24H", coin.MarketCapChangePercentage24H, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#CirculatingSupply", coin.CirculatingSupply, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#TotalSupply", coin.TotalSupply, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#MaxSupply", coin.MaxSupply, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#Ath", coin.Ath, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#AthChangePercentage", coin.AthChangePercentage, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#AthDate", coin.AthDate, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#Atl", coin.Atl, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#AtlChangePercentage", coin.AtlChangePercentage, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#AtlDate", coin.AtlDate, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#Roi", coin.Roi, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	err = client.Set(ctx, coin.ID+"#LastUpdated", coin.LastUpdated, expiry).Err()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
	fmt.Printf("stored: %s\n", coin.ID)
}
