package main

type items struct {
	Stickers []*sticker `json:"stickers"`
	Skins    []*skin    `json:"skins"`
}

type sticker struct {
	MarketHashName string `json:"market_hash_name"`
	Variant        string `json:"variant"`
}

type skin struct {
	MarketHashName string  `json:"market_hash_name"`
	Wear           string  `json:"wear"`
	MinFloat       float64 `json:"min_float"`
	MaxFloat       float64 `json:"max_float"`
}
