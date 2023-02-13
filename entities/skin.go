package entities

type Skin struct {
	MarketHashName string   `json:"market_hash_name"`
	Wear           string   `json:"wear"`
	MinFloat       float64  `json:"min_float"`
	MaxFloat       float64  `json:"max_float"`
	Collection     string   `json:"collection"`
	Containers     []string `json:"containers"`
}

type skinBuilder struct {
	fns []func(*Skin)
}

func (sb *skinBuilder) withMarketHashName(marketHashName string) *skinBuilder {
	sb.fns = append(sb.fns, func(skin *Skin) {
		skin.MarketHashName = marketHashName
	})

	return sb
}
