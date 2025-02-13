package model

import "github.com/shopspring/decimal"

type Chains []Chain

type Chain struct {
	ChainID       string          `json:"chainId"`
	DexID         string          `json:"dexId"`
	URL           string          `json:"url"`
	PairAddress   string          `json:"pairAddress"`
	Labels        []string        `json:"labels"`
	BaseToken     BaseToken       `json:"baseToken"`
	QuoteToken    QuoteToken      `json:"quoteToken"`
	PriceNative   string          `json:"priceNative"`
	PriceUsd      string          `json:"priceUsd"`
	Txns          Txns            `json:"txns"`
	Volume        Volume          `json:"volume"`
	PriceChange   PriceChange     `json:"priceChange"`
	Liquidity     Liquidity       `json:"liquidity"`
	Fdv           decimal.Decimal `json:"fdv"`
	MarketCap     decimal.Decimal `json:"marketCap"`
	PairCreatedAt int             `json:"pairCreatedAt"`
	Info          Info            `json:"info"`
	Boosts        Boosts          `json:"boosts"`
}
type BaseToken struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}
type QuoteToken struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}
type BuySells struct {
	Buys  int `json:"buys"`
	Sells int `json:"sells"`
}
type Txns struct {
	M5 BuySells `json:"m5"`
	H1 BuySells `json:"h1"`
}
type Volume struct {
	M5 decimal.Decimal `json:"m5"`
	H1 decimal.Decimal `json:"h1"`
}
type PriceChange struct {
	AnyAdditionalProperty int `json:"ANY_ADDITIONAL_PROPERTY"`
}
type Liquidity struct {
	Usd   decimal.Decimal `json:"usd"`
	Base  decimal.Decimal `json:"base"`
	Quote decimal.Decimal `json:"quote"`
}
type Websites struct {
	URL string `json:"url"`
}
type Socials struct {
	Platform string `json:"platform"`
	Handle   string `json:"handle"`
}
type Info struct {
	ImageURL string     `json:"imageUrl"`
	Websites []Websites `json:"websites"`
	Socials  []Socials  `json:"socials"`
}
type Boosts struct {
	Active int `json:"active"`
}
