package cmc

// MarketData is a Price and Price Change with json deserialization for USD .
type CurrMarketData struct {
	Data struct {
		S struct {
			Supply float64 `json:"circulating_supply"`
			Quote  struct {
				USD struct {
					Price          float64 `json:"price"`
					Price24hChange float64 `json:"percent_change_24h"`
					MarketCap      float64 `json:"market_cap"`
					Volume         float64 `json:"volume_24h"`
				} `json:"USD"`
			} `json:"quote"`
		} `json:"7653"`
	} `json:"data"`
}

// GetPrice returns the price in USD.
func (md CurrMarketData) GetPrice() float64 {
	return md.Data.S.Quote.USD.Price
}

// GetPriceChange returns the price change during the last 24 hours in percents.
func (md CurrMarketData) GetPriceChange() float64 {
	return md.Data.S.Quote.USD.Price24hChange
}
func (md CurrMarketData) GetMarketCap() float64 {
	return md.Data.S.Quote.USD.MarketCap
}
func (md CurrMarketData) GetVolume() float64 {
	return md.Data.S.Quote.USD.Volume
}
func (md CurrMarketData) GetSupply() float64 {
	return md.Data.S.Supply
}
