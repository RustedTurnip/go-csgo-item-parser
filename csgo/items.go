package csgo

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// itemType is used to categorise what type of item is being dealt with.
type itemType int

const (
	itemTypeUnknown itemType = iota
	itemTypeWeaponGun
	itemTypeWeaponKnife
	itemTypeGloves
	itemTypeCrate
	itemTypeStickerCapsule
)

var (
	// itemIdTypePrefixes stores a number of recognised item Id prefixes
	// if an item is unidentifiable otherwise.
	itemIdTypePrefixes = map[string]itemType{
		"crate_sticker_pack_":   itemTypeStickerCapsule,
		"crate_signature_pack_": itemTypeStickerCapsule,
	}
)

// WeaponQuality represents a skin type, e.g. StatTrak™ or Souvenir
type WeaponQuality string

var (
	qualityNormal   WeaponQuality = "Normal"
	qualityStatTrak WeaponQuality = "StatTrak™"
	qualitySouvenir WeaponQuality = "Souvenir"
)

// itemContainer is just a grouping of relevant items_game items that are parsed
// through getItems.
type itemContainer struct {
	weapons         map[string]*Weapon
	knives          map[string]*Weapon
	gloves          map[string]*Gloves
	crates          map[string]*WeaponCrate
	stickerCapsules map[string]*StickerCapsule
}

// Weapon represents a skinnable item that is also a Weapon in Csgo.
type Weapon struct {
	Id          string
	Name        string
	Description string
}

// mapToWeapon converts the provided map into a Weapon providing
// all required parameters are present and of the correct type.
func mapToWeapon(data map[string]interface{}, prefabs map[string]*itemPrefab, language *language) (*Weapon, error) {

	response := &Weapon{}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("Id (name) missing from Weapon")
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err == nil {
		lang, err := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, "unable to crawl weapon item to path: item_name")
		}

		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		lang, err := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, "unable to crawl weapon item to path: item_description")
		}

		response.Description = lang
	}

	// get info from prefab where missing
	if val, ok := data["prefab"].(string); ok {

		if response.Name == "" {
			response.Name = prefabs[val].name
		}

		if response.Description == "" {
			response.Description = prefabs[val].description
		}
	}

	return response, nil
}

// Gloves represents a special skinnable item that isn't a Weapon.
type Gloves struct {
	Id          string
	Name        string
	Description string
}

// mapToGloves converts the provided map into Gloves providing
// all required parameters are present and of the correct type.
func mapToGloves(data map[string]interface{}, language *language) (*Gloves, error) {

	response := &Gloves{}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl gloves item to path: name")
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl gloves item to path: item_name")
	} else {
		lang, _ := language.lookup(val)
		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl gloves item to path: item_description")
	} else {
		lang, _ := language.lookup(val)
		response.Description = lang
	}

	return response, nil
}

// WeaponCrate represents an openable crate that contains items. The crate's items
// are determined by the linked WeaponSet (item_set).
type WeaponCrate struct {
	Id          string
	Name        string
	Description string

	// WeaponSetId is the ID of the WeaponSet for the item/paintkit combinations
	// available in the crate.
	WeaponSetId string

	// QualityCapability shows whether the crate can produce special skin qualities
	// e.g. Souvenir or StatTrak™
	QualityCapability WeaponQuality
}

// mapToWeaponCrate converts the provided map into a WeaponCrate providing
// all required parameters are present and of the correct type.
func mapToWeaponCrate(data map[string]interface{}, language *language) (*WeaponCrate, error) {

	response := &WeaponCrate{
		QualityCapability: qualityNormal,
	}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl WeaponCrate item to path: name")
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl WeaponCrate item to path: item_name")
	} else {
		lang, err := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to lookup WeaponCrate name (%s) from language", val))
		}

		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		lang, _ := language.lookup(val)
		response.Description = lang
	}

	if val, ok := data["prefab"].(string); ok {

		switch val {
		case "weapon_case":
			response.QualityCapability = qualityStatTrak

		case "weapon_case_souvenirpkg":
			response.QualityCapability = qualitySouvenir
		}
	}

	if val, err := crawlToType[string](data, "tags", "ItemSet", "tag_value"); err == nil {
		response.WeaponSetId = val
	}

	return response, nil

}

// StickerCapsule represents an openable capsule that contains stickers. The capsule's
// stickers are determined by the linked clientLootListId (client_loot_list).
type StickerCapsule struct {
	Id          string
	Name        string
	Description string

	// clientLootListIndex is the index number of the client_loot_list that links the capsule's
	// stickers to the capsule.
	ClientLootListIndex string
}

// mapToStickerCapsule converts the provided map into a StickerCapsule providing
// all required parameters are present and of the correct type.
func mapToStickerCapsule(data map[string]interface{}, language *language) (*StickerCapsule, error) {

	response := &StickerCapsule{}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl StickerCapsule item to path: name")
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err == nil {
		lang, _ := language.lookup(val)
		response.Name = lang
	}

	if val, err := crawlToType[string](data, "tags", "StickerCapsule", "tag_text"); err == nil {
		lang, _ := language.lookup(val)
		response.Name = lang
	}

	if response.Name == "" {
		return nil, errors.New("unable to locate StickerKit's language Name Id")
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		lang, _ := language.lookup(val)
		response.Description = lang
	}

	if val, err := crawlToType[string](data, "attributes", "set supply crate series", "value"); err == nil {
		response.ClientLootListIndex = val
	}

	return response, nil
}

// getItems processes the provided items data and, based on the item's prefab,
// produces the relevant item (e.g. Gloves, Weapon, or crate).
//
// All items are returned within the itemContainer part of the response.
func (c *csgoItems) getItems(prefabs map[string]*itemPrefab) (*itemContainer, error) {

	response := &itemContainer{
		weapons:         make(map[string]*Weapon),
		knives:          make(map[string]*Weapon),
		gloves:          make(map[string]*Gloves),
		crates:          make(map[string]*WeaponCrate),
		stickerCapsules: map[string]*StickerCapsule{},
	}

	items, err := crawlToType[map[string]interface{}](c.items, "items")
	if err != nil {
		return nil, errors.New("items missing from item data") // TODO format error better than this
	}

	for _, itemData := range items {

		itemMap, ok := itemData.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected item format found when fetching items")
		}

		switch getItemType(itemMap, prefabs) {

		case itemTypeWeaponGun:
			w, err := mapToWeapon(itemMap, prefabs, c.language)
			if err != nil {
				return nil, err
			}

			response.weapons[w.Id] = w

		case itemTypeWeaponKnife:
			w, err := mapToWeapon(itemMap, prefabs, c.language)
			if err != nil {
				return nil, err
			}

			response.knives[w.Id] = w

		case itemTypeGloves:
			g, err := mapToGloves(itemMap, c.language)
			if err != nil {
				return nil, err
			}

			response.gloves[g.Id] = g

		case itemTypeCrate:
			c, err := mapToWeaponCrate(itemMap, c.language)
			if err != nil {
				return nil, err
			}

			response.crates[c.Id] = c

		case itemTypeStickerCapsule:
			c, err := mapToStickerCapsule(itemMap, c.language)
			if err != nil {
				return nil, err
			}

			response.stickerCapsules[c.Id] = c
		}
	}

	return response, nil
}

// getItemType attempts to identify an items_game.txt item by assessing its prefab
// (where applicable) or otherwise assessing the contained fields.
func getItemType(data map[string]interface{}, prefabs map[string]*itemPrefab) itemType {

	// attempt to identify from prefab
	prefab, ok := data["prefab"].(string)
	if ok {
		it := getTypeFromPrefab(prefab, prefabs)
		if it != itemTypeUnknown {
			return it
		}
	}

	// attempt to identify from tags
	if val, err := crawlToType[string](data, "tags", "StickerCapsule", "tag_group"); err == nil {
		if val == "StickerCapsule" {
			return itemTypeStickerCapsule
		}
	}

	// attempt to identify from Id prefix
	for prefix, it := range itemIdTypePrefixes {
		if strings.HasPrefix(data["name"].(string), prefix) {
			return it
		}
	}

	return itemTypeUnknown
}
