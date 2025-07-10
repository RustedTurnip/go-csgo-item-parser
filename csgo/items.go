package csgo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type prefabItemConverter func(*csgoItems, int, map[string]interface{}) (interface{}, error)

var (
	// itemPrefabPrefabs provides a mapping of recognised prefab types, to their corresponding
	// item identifying function.
	itemPrefabPrefabs = map[string]prefabItemConverter{
		"primary": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToWeapon(index, data, items.prefabs, items.language)
		},

		"secondary": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToWeapon(index, data, items.prefabs, items.language)
		},

		"melee_unusual": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToWeapon(index, data, items.prefabs, items.language)
		},

		"hands": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToGloves(index, data, items.language)
		},

		"equipment": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			pfID, _ := data["prefab"].(string)

			// some equipment items are weapons, checking for item gear slot
			// differentiates between them
			if pf, _ := items.prefabs[pfID]; pf.itemGearSlot == "melee" {
				return mapToWeapon(index, data, items.prefabs, items.language)
			}

			return mapToEquipment(index, data, items.language)
		},

		"weapon_case": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToWeaponCrate(index, data, items.language)
		},

		"weapon_case_souvenirpkg": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToWeaponCrate(index, data, items.language)
		},

		"weapon_case_base": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			// weapon crate cast
			if _, err := crawlToType[string](data, "tags", "ItemSet", "tag_value"); err == nil {
				return mapToWeaponCrate(index, data, items.language)
			}

			// if it is a set (identified through revolving_loot_lists)
			if val, err := crawlToType[string](data, "attributes", "set supply crate series", "value"); err == nil {
				if clientLootListId, ok := items.revolvingLootLists[val]; ok {
					itemType, listItems := crawlClientLootLists(clientLootListId, items.clientLootLists)

					switch itemType {
					case clientLootListItemTypeSticker:
						return mapToStickerCapsule(index, data, listItems, items.language)
					}
				}
			}

			// ignore items where the contained list is located through the key "loot_list_name" as these
			// are either not capsules (but instead the StoreItem representing them), or are the duplicate
			// self-opening version of a set e.g. selfopeningitem_crate_sticker_pack_riptide_surfshop

			// Can be additionally split into:
			// - Operator Dossier - (TODO)
			// - Music Kit capsule (TODO)
			// - Collectibles Collections (TODO)

			return nil, nil
		},

		"csgo_tool": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToTool(index, data, items.language)
		},
		"customplayertradable": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToCharacter(index, data, items.language)
		},

		"collectible": func(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
			return mapToCollectible(index, data, items.language)
		},

		// these seem like placeholder values, they break collectables
		"collectible_untradable_coin": func(ci *csgoItems, i int, m map[string]interface{}) (interface{}, error) {
			return nil, nil
		},
	}
)

// WeaponQuality represents a skin type, e.g. StatTrak™ or Souvenir
type WeaponQuality string

var (
	QualityNormal   WeaponQuality = ""
	QualityStatTrak WeaponQuality = "StatTrak™"
	QualitySouvenir WeaponQuality = "Souvenir"
)

// itemContainer is just a grouping of relevant items_game items that are parsed
// through getItems.
type itemContainer struct {
	weapons         map[string]*Weapon
	knives          map[string]*Weapon
	gloves          map[string]*Gloves
	equipment       map[string]*Equipment
	crates          map[string]*WeaponCrate
	stickerCapsules map[string]*StickerCapsule
	tools           map[string]*Tool
	characters      map[string]*Character
	collectables    map[string]*Collectible
}

// Weapon represents a skinnable item that is also a Weapon in Csgo.
type Weapon struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// mapToWeapon converts the provided map into a Weapon providing
// all required parameters are present and of the correct type.
func mapToWeapon(index int, data map[string]interface{}, prefabs map[string]*itemPrefab, language *language) (*Weapon, error) {
	response := &Weapon{
		Index: index,
	}

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
		spl := strings.Split(val, " ")

		if response.Name == "" {
			response.Name = prefabs[spl[0]].name
		}

		if response.Description == "" {
			response.Description = prefabs[spl[0]].description
		}
	}

	return response, nil
}

// Equipment represents miscellaneous items in game that don't
// constitute weapons.
type Equipment struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// mapToWeapon converts the provided map into a Weapon providing
// all required parameters are present and of the correct type.
func mapToEquipment(index int, data map[string]interface{}, language *language) (*Equipment, error) {
	response := &Equipment{
		Index: index,
	}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("Id (name) missing from Equipment")
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err == nil {
		lang, err := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, "unable to crawl equipment item to path: item_name")
		}

		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		lang, err := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, "unable to crawl equipment item to path: item_description")
		}

		response.Description = lang
	}

	return response, nil
}

// Gloves represents a special skinnable item that isn't a Weapon.
type Gloves struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// mapToGloves converts the provided map into Gloves providing
// all required parameters are present and of the correct type.
func mapToGloves(index int, data map[string]interface{}, language *language) (*Gloves, error) {
	response := &Gloves{
		Index: index,
	}

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
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// WeaponSetId is the ID of the WeaponSet for the item/Paintkit combinations
	// available in the crate.
	WeaponSetIds []string `json:"weaponSetIds"`

	// QualityCapability shows whether the crate can produce special skin qualities
	// e.g. Souvenir or StatTrak™
	QualityCapability WeaponQuality `json:"qualityCapability"`
}

// mapToWeaponCrate converts the provided map into a WeaponCrate providing
// all required parameters are present and of the correct type.
func mapToWeaponCrate(index int, data map[string]interface{}, language *language) (*WeaponCrate, error) {
	response := &WeaponCrate{
		Index:             index,
		QualityCapability: QualityNormal,
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
			response.QualityCapability = QualityStatTrak

		case "weapon_case_souvenirpkg":
			response.QualityCapability = QualitySouvenir
		}
	}

	if val, err := crawlToType[string](data, "tags", "ItemSet", "tag_value"); err == nil {
		response.WeaponSetIds = append(response.WeaponSetIds, val)
	}

	if response.WeaponSetIds == nil {
		// the earliest crates were comprised of these sets, but the link doesn't exist within the
		// items_game file.
		response.WeaponSetIds = []string{
			"set_lake",
			"set_italy",
			"set_safehouse",
		}
	}

	return response, nil

}

// StickerCapsule represents an openable capsule that contains stickers. The capsule's
// stickers are determined by the linked clientLootListId (client_loot_list).
type StickerCapsule struct {
	Id          string   `json:"id"`
	Index       int      `json:"index"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	StickerKits []string `json:"stickerKits"`
}

// mapToStickerCapsule converts the provided map into a StickerCapsule providing
// all required parameters are present and of the correct type.
func mapToStickerCapsule(index int, data map[string]interface{}, stickers []string, language *language) (*StickerCapsule, error) {
	response := &StickerCapsule{
		Index:       index,
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

// Tool represents consumable inventory only items
type Tool struct {
	Id    string `json:"id"`
	Index int    `json:"index"`
	Name  string `json:"name"`
}

// mapToTool converts the provided map into a Tool providing
// all required parameters are present and of the correct type.
func mapToTool(index int, data map[string]interface{}, language *language) (*Tool, error) {
	response := &Tool{
		Index: index,
	}
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl Tool item to path: name")
	} else {
		response.Id = val
	}

	if val, err := crawlToType[string](data, "item_name"); err == nil {
		lang, _ := language.lookup(val)
		if lang == "" {
			lang, _ = language.lookup(val[1:])
		}
		response.Name = lang
	} else if val, err := crawlToType[string](data, "item_type_name"); err == nil {
		lang, _ := language.lookup(val)
		if lang == "" {
			lang, _ = language.lookup(val[1:])
		}
		response.Name = lang
	}

	if response.Name == "" {
		return nil, errors.New("unable to locate Tools's language Name Id" + fmt.Sprintf("%+v", response))
	}

	return response, nil
}

// Character represents a skin that a player can use ingame
type Character struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RarityId    string `json:"rarityId"`
}

// mapToCharacter converts the provided map into a Character providing
// all required parameters are present and of the correct type.
func mapToCharacter(index int, data map[string]interface{}, language *language) (*Character, error) {
	response := &Character{
		Index: index,
	}
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl Tool item to path: name")
	} else {
		response.Id = val
	}

	if val, err := crawlToType[string](data, "item_name"); err == nil {
		lang, _ := language.lookup(val[1:])
		response.Name = lang
	}

	if val, err := crawlToType[string](data, "item_description"); err == nil {
		lang, _ := language.lookup(val[1:])
		response.Description = lang
	}

	if val, err := crawlToType[string](data, "item_rarity"); err == nil {
		response.RarityId = val
	}

	if response.Name == "" {
		return nil, errors.New("unable to locate Character's language Name Id" + fmt.Sprintf("%+v", response))
	}
	if response.Description == "" {
		return nil, errors.New("unable to locate Character's language Description Id" + fmt.Sprintf("%+v", response))
	}
	if response.RarityId == "" {
		return nil, errors.New("unable to locate Character's language Rarity Id" + fmt.Sprintf("%+v", response))
	}

	return response, nil
}

type Collectible struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// mapToCharacter converts the provided map into a Character providing
// all required parameters are present and of the correct type.
func mapToCollectible(index int, data map[string]interface{}, language *language) (*Collectible, error) {
	response := &Collectible{
		Index: index,
	}
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "unable to crawl Collectible item to path: name")
	} else {
		response.Id = val
	}

	if val, err := crawlToType[string](data, "item_name"); err == nil {
		lang, _ := language.lookup(val[1:])
		response.Name = lang
	}

	if val, err := crawlToType[string](data, "item_description"); err == nil {
		lang, _ := language.lookup(val[1:])
		response.Description = lang
	}

	if response.Name == "" {
		return nil, errors.New("unable to locate Collectible's language Name Id" + fmt.Sprintf("%+v", response))
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
		equipment:       make(map[string]*Equipment),
		crates:          make(map[string]*WeaponCrate),
		stickerCapsules: make(map[string]*StickerCapsule),
		tools:           make(map[string]*Tool),
		characters:      make(map[string]*Character),
		collectables:    make(map[string]*Collectible),
	}

	items, err := crawlToType[map[string]interface{}](c.items, "items")
	if err != nil {
		return nil, errors.Wrap(err, "items (at path \"items\") missing from item data")
	}

	for index, itemData := range items {
		if index == "default" {
			continue
		}

		iIndex, err := strconv.Atoi(index)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret item index (%s) as int", iIndex))
		}

		itemMap, ok := itemData.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected item format found when fetching items")
		}

		converted, err := convertItem(c, iIndex, itemMap)
		if err != nil {
			return nil, err
		}

		switch t := converted.(type) {
		case *Weapon:
			if itemMap["prefab"].(string) == "melee_unusual" {
				response.knives[t.Id] = t
				continue
			}

			response.weapons[t.Id] = t

		case *Gloves:
			response.gloves[t.Id] = t

		case *Equipment:
			response.equipment[t.Id] = t

		case *WeaponCrate:
			response.crates[t.Id] = t

		case *StickerCapsule:
			response.stickerCapsules[t.Id] = t

		case *Tool:
			response.tools[t.Id] = t

		case *Character:
			response.characters[t.Id] = t

		case *Collectible:
			response.collectables[t.Id] = t
		}
	}

	return response, nil
}

// getItemType attempts to identify an items_game.txt item by assessing its prefab
// (where applicable) or otherwise assessing the contained fields.
func convertItem(items *csgoItems, index int, data map[string]interface{}) (interface{}, error) {
	prefab, ok := data["prefab"].(string)
	if !ok {
		return nil, nil
	}

	converter := getPrefabConversionFunc(prefab, items.prefabs)
	if converter == nil {
		return nil, nil
	}

	return converter(items, index, data)
}

// getPrefabConversionFunc attempts to identify the correct conversion function for the item data map
// from the item's prefab.
func getPrefabConversionFunc(prefabId string, prefabs map[string]*itemPrefab) prefabItemConverter {
	if prefabId == "collectible_untradable_coin" {
		return nil
	}
	if converter, ok := itemPrefabPrefabs[prefabId]; ok {
		return converter
	}

	prefab, ok := prefabs[prefabId]

	// if prefab isn't recognised
	if !ok {
		return nil
	}

	for _, parent := range prefab.parentPrefabs {
		// if prefab item type is unknown, but prefab has parent, crawl further
		c := getPrefabConversionFunc(parent, prefabs)
		if c == nil {
			// if not found with this parent, move onto the next parent (prefabs can have multiple parents)
			continue
		}

		return c
	}

	// at this point, we have gotten to the roots of each parent prefab without finding anything
	// so return with an unknown response
	return nil
}
