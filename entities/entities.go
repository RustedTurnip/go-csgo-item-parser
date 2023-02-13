package entities

type Items struct {
	Stickers []*Sticker `json:"stickers"`
	Skins    []*Skin    `json:"skins"`
}
