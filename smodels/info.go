package smodels

type Info struct {
	Height      uint64  `json:"height"`
	TopEscrow   float64 `json:"top_escrow"`
	Price       float64 `json:"price"`
	PriceChange float64 `json:"price_24h_change"`
	MarketCap   float64 `json:"market_cap"`
	Volume      float64 `json:"volume_24h"`
	Supply      float64 `json:"circulating_supply"`
}
