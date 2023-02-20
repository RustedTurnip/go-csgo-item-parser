package csgo

import (
	"errors"
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

// qualityCapability represents a skin type, e.g. StatTrak™ or Souvenir
type qualityCapability string

// skinnableItem represents any item that can be represented as a skin with a
// Market Hash Name. Internally, to derive the Market Hash Name, we require a
// descriptive Name Id for the language file, and whether the item is special.
type skinnableItem interface {
	getLanguageNameId() string
	getSpecial() bool
}

var (
	qualityNormal   qualityCapability = "Normal"
	qualityStatTrak qualityCapability = "StatTrak™"
	qualitySouvenir qualityCapability = "Souvenir"
)

// itemContainer is just a grouping of relevant items_game items that are parsed
// through getItems.
type itemContainer struct {
	weapons         map[string]*weapon
	knives          map[string]*weapon
	gloves          map[string]*gloves
	crates          map[string]*weaponCrate
	stickerCapsules map[string]*stickerCapsule
}

// weapon represents a skinnable item that is also a weapon in Csgo.
type weapon struct {
	id          string
	name        string
	description string
}

// mapToWeapon converts the provided map into a weapon providing
// all required parameters are present and of the correct type.
func mapToWeapon(data map[string]interface{}, prefabs map[string]*itemPrefab, language *language) (*weapon, error) {

	response := &weapon{}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("Id (name) missing from weapon")
	} else {
		response.id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err == nil {
		lang, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		lang, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.description = lang
	}

	// get info from prefab where missing
	if val, ok := data["prefab"].(string); ok {

		if response.name == "" {
			response.name = prefabs[val].name
		}

		if response.description == "" {
			response.description = prefabs[val].description
		}
	}

	return response, nil
}

// gloves represents a special skinnable item that isn't a weapon.
type gloves struct {
	Id          string
	Name        string
	Description string
}

// mapToGloves converts the provided map into gloves providing
// all required parameters are present and of the correct type.
func mapToGloves(data map[string]interface{}, language *language) (*gloves, error) {

	response := &gloves{}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("Id (name) missing from weapon") // TODO improve error
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.New("language Name Id (item_name) missing from weapon") // TODO improve error
	} else {
		lang, _ := language.lookup(val)
		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err != nil {
		return nil, errors.New("language Description Id (item_description) missing from weapon") // TODO improve error
	} else {
		lang, _ := language.lookup(val)
		response.Description = lang
	}

	return response, nil
}

// weaponCrate represents an openable crate that contains items. The crate's items
// are determined by the linked weaponSet (item_set).
type weaponCrate struct {
	id                    string
	languageNameId        string
	languageDescriptionId string

	// collectionId is the ID of the weaponSet for the item/paintkit combinations
	// available in the crate.
	collectionId string

	// qualityCapability shows whether the crate can produce special skin qualities
	// e.g. Souvenir or StatTrak™
	qualityCapability qualityCapability
}

func (c *weaponCrate) getLanguageNameId() string {
	return c.languageNameId
}

// mapToWeaponCrate converts the provided map into a weaponCrate providing
// all required parameters are present and of the correct type.
func mapToWeaponCrate(data map[string]interface{}) (*weaponCrate, error) {

	response := &weaponCrate{
		qualityCapability: qualityNormal,
	}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("Id (Name) missing from weaponCrate") // TODO improve error
	} else {
		response.id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.New("language Name Id (item_name) missing from weaponCrate") // TODO improve error
	} else {
		response.languageNameId = val
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		response.languageDescriptionId = val
	}

	if val, ok := data["prefab"].(string); ok {

		switch val {
		case "weapon_case":
			response.qualityCapability = qualityStatTrak

		case "weapon_case_souvenirpkg":
			response.qualityCapability = qualitySouvenir
		}
	}

	if val, err := crawlToType[string](data, "tags", "ItemSet", "tag_value"); err == nil {
		response.collectionId = val
	}

	return response, nil

}

// stickerCapsule represents an openable capsule that contains stickers. The capsule's
// stickers are determined by the linked clientLootListId (client_loot_list).
type stickerCapsule struct {
	id                    string
	languageNameId        string
	languageDescriptionId string

	// clientLootListIndex is the index number of the client_loot_list that links the capsule's
	// stickers to the capsule.
	clientLootListIndex string
}

func (c *stickerCapsule) getLanguageNameId() string {
	return c.languageNameId
}

// mapToStickerCapsule converts the provided map into a stickerCapsule providing
// all required parameters are present and of the correct type.
func mapToStickerCapsule(data map[string]interface{}) (*stickerCapsule, error) {

	response := &stickerCapsule{}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("Id (name) missing from stickerCapsule") // TODO improve error
	} else {
		response.id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err == nil {
		response.languageNameId = val
	}

	if val, err := crawlToType[string](data, "tags", "StickerCapsule", "tag_text"); err == nil {
		response.languageNameId = val
	}

	if response.languageNameId == "" {
		return nil, errors.New("unable to locate StickerKit's language Name Id")
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		response.languageDescriptionId = val
	}

	if val, err := crawlToType[string](data, "attributes", "set supply crate series", "value"); err == nil {
		response.clientLootListIndex = val
	}

	return response, nil
}

// getItems processes the provided items data and, based on the item's prefab,
// produces the relevant item (e.g. gloves, weapon, or crate).
//
// All items are returned within the itemContainer part of the response.
func (c *csgoItems) getItems(prefabs map[string]*itemPrefab) (*itemContainer, error) {

	response := &itemContainer{
		weapons:         make(map[string]*weapon),
		knives:          make(map[string]*weapon),
		gloves:          make(map[string]*gloves),
		crates:          make(map[string]*weaponCrate),
		stickerCapsules: map[string]*stickerCapsule{},
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

			response.weapons[w.id] = w

		case itemTypeWeaponKnife:
			w, err := mapToWeapon(itemMap, prefabs, c.language)
			if err != nil {
				return nil, err
			}

			response.knives[w.id] = w

		case itemTypeGloves:
			g, err := mapToGloves(itemMap, c.language)
			if err != nil {
				return nil, err
			}

			response.gloves[g.Id] = g

		case itemTypeCrate:
			c, err := mapToWeaponCrate(itemMap)
			if err != nil {
				return nil, err
			}

			response.crates[c.id] = c

		case itemTypeStickerCapsule:
			c, err := mapToStickerCapsule(itemMap)
			if err != nil {
				return nil, err
			}

			response.stickerCapsules[c.id] = c
		}
	}

	return response, nil
}

// TODO comment
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
