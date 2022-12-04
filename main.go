package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	logger    = log.New()
	ctx       = context.Background()
	frequency *int
	pages     *int
	currency  *string
	order     *string
	expiry    *int
	hostname  *string
	password  *string
	db        *int
	metrics   *string
	wg        sync.WaitGroup
	tail      []string
	cgHits    = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cg_hit",
			Help: "Number of times the cache got data",
		},
	)
	cgMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cg_miss",
			Help: "Number of times the cache missed data",
		},
	)
)

const (
	CoinGeckoMarkets = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&order=%s&page=%d"
	CoinGeckoCoin    = "https://api.coingecko.com/api/v3/coins/%s"
)

func init() {
	frequency = flag.Int("frequency", 1, "seconds between updates")
	pages = flag.Int("pages", 1, "number of pages (100 coin each) to pull from")
	currency = flag.String("currency", "usd", "currency to use")
	order = flag.String("order", "market_cap_desc", "sort key: market_cap_desc, gecko_desc, gecko_asc, market_cap_asc, market_cap_desc, volume_asc, volume_desc, id_asc, id_desc")
	expiry = flag.Int("expiry", 60, "number of seconds to keep entries in the cache")
	hostname = flag.String("hostname", "localhost:6379", "connection address for redis")
	password = flag.String("password", "", "redis password")
	db = flag.Int("db", 0, "redis db to use")
	metrics = flag.String("metrics", ":6380", "port for metrics server")
	flag.Parse()
	tail = flag.Args()
	logger.Out = os.Stdout
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     *hostname,
		Password: *password,
		DB:       *db,
	})

	go gather(rdb)

	prometheus.MustRegister(cgHits)
	prometheus.MustRegister(cgMisses)

	http.Handle("/metrics", promhttp.Handler())
	logger.Error(http.ListenAndServe(*metrics, nil))
}

func gather(rdb *redis.Client) {
	pager := 1

	for {

		coinsData, err := GetMarketData(pager)
		if err != nil {
			logger.Error(err)
			time.Sleep(time.Duration(*frequency) * time.Second)
			continue
		}

		for _, coin := range coinsData {
			wg.Add(1)
			go Store(&wg, rdb, coin, time.Duration(*expiry)*time.Second)
		}
		fmt.Println("waiting for market storage")
		wg.Wait()

		if pager == 1 {
			for _, coin := range tail {

				coinsData, err := GetCoinData(coin)
				if err != nil {
					logger.Error(err)
					continue
				}
				wg.Add(1)
				go Store(&wg, rdb, coinsData, time.Duration(*expiry)*time.Second)
			}
			fmt.Println("waiting for specific storage")
			wg.Wait()
		}

		pager++
		if pager > *pages {
			pager = 1
			time.Sleep(time.Duration(*frequency) * time.Second)
		}
	}
}

// GetMarketData retrieves the array of current prices from coingecko
func GetMarketData(page int) ([]MarketInfo, error) {
	var prices []MarketInfo

	req, err := http.NewRequest("GET", fmt.Sprintf(CoinGeckoMarkets, *currency, *order, page), nil)
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
	fmt.Printf("Retrived page %d\n", page)

	if resp.StatusCode == 429 {
		cgMisses.Inc()
		return prices, errors.New("being rate limited by coingecko")
	}

	results, err := io.ReadAll(resp.Body)
	if err != nil {
		return prices, err
	}

	err = json.Unmarshal(results, &prices)
	if err != nil {
		return prices, err
	}

	cgHits.Inc()
	return prices, nil
}

// GetCoinData retrive data on a single id
func GetCoinData(id string) (MarketInfo, error) {
	var price MarketInfo
	var coinPrice CoinData

	req, err := http.NewRequest("GET", fmt.Sprintf(CoinGeckoCoin, id), nil)
	if err != nil {
		return price, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("accept", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return price, err
	}
	fmt.Printf("Retrived id %s\n", id)

	if resp.StatusCode == 429 {
		cgMisses.Inc()
		return price, errors.New("being rate limited by coingecko")
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return price, err
	}

	err = json.Unmarshal(result, &coinPrice)
	if err != nil {
		return price, err
	}

	price = MarketInfo{
		ID:                           id,
		Symbol:                       coinPrice.Symbol,
		Name:                         coinPrice.Name,
		Image:                        coinPrice.Image.Thumb,
		CurrentPrice:                 float64(coinPrice.MarketData.CurrentPrice.Usd),
		MarketCap:                    float64(coinPrice.MarketData.MarketCap.Usd),
		MarketCapRank:                coinPrice.MarketCapRank,
		FullyDilutedValuation:        float64(coinPrice.MarketData.FullyDilutedValuation.Usd),
		TotalVolume:                  float64(coinPrice.MarketData.TotalVolume.Usd),
		High24H:                      float64(coinPrice.MarketData.High24H.Usd),
		Low24H:                       float64(coinPrice.MarketData.Low24H.Usd),
		PriceChange24H:               coinPrice.MarketData.PriceChange24H,
		PriceChangePercentage24H:     coinPrice.MarketData.PriceChangePercentage24H,
		MarketCapChange24H:           coinPrice.MarketData.MarketCapChange24H,
		MarketCapChangePercentage24H: coinPrice.MarketData.MarketCapChangePercentage24H,
		CirculatingSupply:            coinPrice.MarketData.CirculatingSupply,
		TotalSupply:                  coinPrice.MarketData.TotalSupply,
		MaxSupply:                    coinPrice.MarketData.MaxSupply,
		Ath:                          float64(coinPrice.MarketData.Ath.Usd),
		AthChangePercentage:          coinPrice.MarketData.AthChangePercentage.Usd,
		AthDate:                      coinPrice.MarketData.AthDate.Usd,
		Atl:                          float64(coinPrice.MarketData.Ath.Usd),
		AtlChangePercentage:          coinPrice.MarketData.AtlChangePercentage.Usd,
		AtlDate:                      coinPrice.MarketData.AtlDate.Usd,
		LastUpdated:                  coinPrice.LastUpdated,
	}

	cgHits.Inc()
	return price, nil
}

// Store puts the coins values into redis
func Store(wg *sync.WaitGroup, client *redis.Client, coin MarketInfo, expiry time.Duration) {
	defer wg.Done()

	err := client.Set(ctx, coin.ID+"#Symbol", coin.Symbol, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#Name", coin.Name, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#Image", coin.Image, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#CurrentPrice", coin.CurrentPrice, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#MarketCap", coin.MarketCap, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#MarketCapRank", coin.MarketCapRank, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#FullyDilutedValuation", coin.FullyDilutedValuation, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#TotalVolume", coin.TotalVolume, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#High24H", coin.High24H, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#Low24H", coin.Low24H, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#PriceChange24H", coin.PriceChange24H, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#PriceChangePercentage24H", coin.PriceChangePercentage24H, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#MarketCapChange24H", coin.MarketCapChange24H, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#MarketCapChangePercentage24H", coin.MarketCapChangePercentage24H, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#CirculatingSupply", coin.CirculatingSupply, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#TotalSupply", coin.TotalSupply, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#MaxSupply", coin.MaxSupply, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#Ath", coin.Ath, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#AthChangePercentage", coin.AthChangePercentage, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#AthDate", coin.AthDate, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#Atl", coin.Atl, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#AtlChangePercentage", coin.AtlChangePercentage, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#AtlDate", coin.AtlDate, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	err = client.Set(ctx, coin.ID+"#LastUpdated", coin.LastUpdated, expiry).Err()
	if err != nil {
		logger.Error(err)
	}
	fmt.Printf("stored: %s\n", coin.ID)
}
