package cmc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	coingecko "github.com/superoo7/go-gecko/v3"
)

const (
	tezosPriceURL = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=%s&ids=oasis_network&order=market_cap_desc&per_page=100&page=1&sparkline=false&price_change_percentage=24h"
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
	GetTezosMarketData(curr string) (md MarketInfo, err error)
}

var AvailableCurrencies = map[string]bool{"usd": true, "eur": true, "gbp": true, "cny": true}

// CoinGecko is market data provider.
type CoinGecko struct {
	Cache *cache.Cache
}

func NewCoinGecko() *CoinGecko {
	return &CoinGecko{cache.New(cacheTTL, cacheTTL)}
}

// GetTezosMarketData gets the tezos prices and price change from CoinGecko API.
func (c CoinGecko) GetTezosMarketData(curr string) (md MarketInfo, err error) {
	if !AvailableCurrencies[curr] {
		return md, fmt.Errorf("Not available currency: %s", curr)
	}

	md, err = c.GetTezosMarketDataByCurr(curr)
	if err != nil {
		return md, err
	}

	return md, nil
}

func (c CoinGecko) GetTezosMarketDataByCurr(curr string) (md CurrMarketData, err error) {
	cacheKey := fmt.Sprintf(marketInfoKey, curr)
	if marketData, isFound := c.Cache.Get(cacheKey); isFound {
		return marketData.(CurrMarketData), nil
	}

	cg := coingecko.NewClient(nil)
	b, err := cg.MakeReq(fmt.Sprintf(tezosPriceURL, curr))
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
	c.Cache.Add(cacheKey, tmd[0], cacheTTL)

	return tmd[0], nil
}
