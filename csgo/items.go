package csgo

import (
	"fmt"
	"github.com/pkg/errors"
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

type prefabItemConverter func(*csgoItems, map[string]interface{}) (interface{}, error)

var (
	// itemPrefabPrefabs is a map of all prefabs that exist against item prefabs
	// that we need to track.
	itemPrefabPrefabs = map[string]itemType{

		"": itemTypeUnknown,

		// Guns
		"primary":       itemTypeWeaponGun,   // covers ay Weapon that can be primary (e.g. smg)
		"secondary":     itemTypeWeaponGun,   // covers any Weapon that can be secondary (e.g. pistol)
		"melee_unusual": itemTypeWeaponKnife, // covers all tradable Knives

		// Gloves
		"hands": itemTypeGloves, // covers Gloves

		// crates
		"weapon_case":             itemTypeCrate,
		"weapon_case_souvenirpkg": itemTypeCrate,

		// stickers
		"sticker_capsule": itemTypeStickerCapsule,
	}

	itemPrefabPrefabs2 = map[string]prefabItemConverter{

		"primary": func(items *csgoItems, data map[string]interface{}) (interface{}, error) {
			return mapToWeapon(data, items.prefabs, items.language)
		},

		"secondary": func(items *csgoItems, data map[string]interface{}) (interface{}, error) {
			return mapToWeapon(data, items.prefabs, items.language)
		},

		"melee_unusual": func(items *csgoItems, data map[string]interface{}) (interface{}, error) {
			return mapToWeapon(data, items.prefabs, items.language)
		},

		"hands": func(items *csgoItems, data map[string]interface{}) (interface{}, error) {
			return mapToGloves(data, items.language)
		},

		"weapon_case": func(items *csgoItems, data map[string]interface{}) (interface{}, error) {
			return mapToWeaponCrate(data, items.language)
		},

		"weapon_case_souvenirpkg": func(items *csgoItems, data map[string]interface{}) (interface{}, error) {
			return mapToWeaponCrate(data, items.language)
		},

		"weapon_case_base": func(items *csgoItems, data map[string]interface{}) (interface{}, error) {

			// weapon crate cast
			if _, err := crawlToType[string](data, "tags", "ItemSet", "tag_value"); err == nil {
				return mapToWeaponCrate(data, items.language)
			}

			// if it is a set (identified through revolving_loot_lists)
			if val, err := crawlToType[string](data, "attributes", "set supply crate series", "value"); err == nil {
				if clientLootListId, ok := items.revolvingLootLists[val]; ok {
					itemType, listItems := crawlClientLootLists(clientLootListId, items.clientLootLists)

					switch itemType {
					case clientLootListItemTypeSticker:
						return mapToStickerCapsule(data, listItems, items.language)
					}
				}
			}

			// if it is a set (identified through client_loot_lists)
			if val, err := crawlToType[string](data, "loot_list_name"); err == nil {
				itemType, listItems := crawlClientLootLists(val, items.clientLootLists)

				switch itemType {
				case clientLootListItemTypeSticker:
					return mapToStickerCapsule(data, listItems, items.language)
				}
			}

			// Can be split into:
			// - Sticker Pack (Capsule) - (Deduce from revolving_loot_lists)
			// - Operator Dossier - (Deduce from revolving_loot_lists) (TODO add when we support characters)
			// - Music Kit capsule (TODO add when we support music kits)

			return nil, nil // TODO
		},
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
	StickerKits []string
}

// mapToStickerCapsule converts the provided map into a StickerCapsule providing
// all required parameters are present and of the correct type.
func mapToStickerCapsule(data map[string]interface{}, stickers []string, language *language) (*StickerCapsule, error) {

	response := &StickerCapsule{
		StickerKits: stickers,
	}

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

	return response, nil
}

// getItems processes the provided items data and, based on the item's prefab,
// produces the relevant item (e.g. Gloves, Weapon, or crate).
//
// All items are returned within the itemContainer part of the response.
func (c *csgoItems) getItems() (*itemContainer, error) {

	response := &itemContainer{
		weapons:         make(map[string]*Weapon),
		knives:          make(map[string]*Weapon),
		gloves:          make(map[string]*Gloves),
		crates:          make(map[string]*WeaponCrate),
		stickerCapsules: make(map[string]*StickerCapsule),
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

		converted, err := convertItem(c, itemMap)
		if err != nil {
			return nil, err
		}

		switch converted.(type) {

		case *Weapon:
			if itemMap["prefab"].(string) == "melee_unusual" {
				response.knives[converted.(*Weapon).Id] = converted.(*Weapon)
				continue
			}

			response.weapons[converted.(*Weapon).Id] = converted.(*Weapon)

		case *Gloves:
			response.gloves[converted.(*Gloves).Id] = converted.(*Gloves)

		case *WeaponCrate:
			response.crates[converted.(*WeaponCrate).Id] = converted.(*WeaponCrate)

		case *StickerCapsule:
			response.stickerCapsules[converted.(*StickerCapsule).Id] = converted.(*StickerCapsule)

		}
	}

	return response, nil
}

// getItemType attempts to identify an items_game.txt item by assessing its prefab
// (where applicable) or otherwise assessing the contained fields.
func convertItem(items *csgoItems, data map[string]interface{}) (interface{}, error) {

	prefab, ok := data["prefab"].(string)
	if !ok {
		return nil, nil
	}

	converter := getPrefabConversionFunc(prefab, items.prefabs)
	if converter == nil {
		return nil, nil
	}

	return converter(items, data)
}

// TODO comment
func getPrefabConversionFunc(prefabId string, prefabs map[string]*itemPrefab) prefabItemConverter {

	if converter, ok := itemPrefabPrefabs2[prefabId]; ok {
		return converter
	}

	prefab, ok := prefabs[prefabId]

	// if prefab isn't recognised
	if !ok {
		return nil
	}

	// if at root of prefab tree, return unknown
	if prefab.parentPrefab == "" {
		return nil
	}

	// if prefab item type is unknown, but prefab has parent, crawl further
	return getPrefabConversionFunc(prefab.parentPrefab, prefabs)
}
