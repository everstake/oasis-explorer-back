package cmc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	coingecko "github.com/superoo7/go-gecko/v3"
)

const (
	oasisPriceURL = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&ids=oasis-network&order=market_cap_desc&per_page=100&page=1&sparkline=false&price_change_percentage=24h"
	oasisInfoURL  = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?id=7653"
	cacheTTL      = 1 * time.Minute
	marketInfoKey = "market_info_%s"
)

// MarketInfo is the interface getting prices and price changes.
type MarketInfo interface {
	GetPrice() float64
	GetPriceChange() float64
	GetMarketCap() float64
	GetVolume() float64
	GetSupply() float64
}

// MarketDataProvider is an interface for getting actual price and price changes.
type MarketDataProvider interface {
	GetOasisMarketData(curr, key string) (md MarketInfo, err error)
}

var AvailableCurrencies = map[string]bool{"usd": true, "eur": true, "gbp": true, "cny": true}

// CoinGecko is market data provider.
type CoinGecko struct {
	Cache *cache.Cache
}

func NewCoinGecko() *CoinGecko {
	return &CoinGecko{cache.New(cacheTTL, cacheTTL)}
}

// GetOasisMarketData gets the oasis prices and price change from CoinGecko API.
func (c CoinGecko) GetOasisMarketData(curr, key string) (md MarketInfo, err error) {
	if !AvailableCurrencies[curr] {
		return md, fmt.Errorf("Not available currency: %s", curr)
	}

	md, err = c.GetOasisMarketDataByCurrCMC(curr, key)
	if err != nil {
		return md, err
	}

	return md, nil
}

func (c CoinGecko) GetOasisMarketDataByCurr(curr string) (md CurrMarketData, err error) {
	cacheKey := fmt.Sprintf(marketInfoKey, curr)
	if marketData, isFound := c.Cache.Get(cacheKey); isFound {
		return marketData.(CurrMarketData), nil
	}

	cg := coingecko.NewClient(nil)
	b, err := cg.MakeReq(fmt.Sprintf(oasisPriceURL, curr))
	if err != nil {
		return md, err
	}
	var tmd []CurrMarketData
	err = json.Unmarshal(b, &tmd)
	if err != nil {
		return md, err
	}
	if len(tmd) != 1 {
		return md, fmt.Errorf("got enexpected number of entries")
	}

	//Save into cache error can be skipped
	err = c.Cache.Add(cacheKey, tmd[0], cacheTTL)
	if err != nil {
		return md, err
	}

	return tmd[0], nil
}

func (c CoinGecko) GetOasisMarketDataByCurrCMC(curr string, key string) (md CurrMarketData, err error) {
	cacheKey := fmt.Sprintf(marketInfoKey, curr)
	if marketData, isFound := c.Cache.Get(cacheKey); isFound {
		return marketData.(CurrMarketData), nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", oasisInfoURL, nil)
	if err != nil {
		return md, fmt.Errorf("cannot make a request")
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", key)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return md, err
	}
	err = json.Unmarshal(body, &md)
	if err != nil {
		return md, err
	}

	//Save into cache error can be skipped
	err = c.Cache.Add(cacheKey, md, cacheTTL)
	if err != nil {
		return md, err
	}

	return md, nil
}
