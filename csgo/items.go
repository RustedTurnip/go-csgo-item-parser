package csgo

import (
	"errors"
)

// itemType is used to categorise what type of item is being dealt with.
type itemType int

const (
	itemTypeUnknown itemType = iota
	itemTypeWeapon
	itemTypeGloves
	itemTypeCrate
)

// qualityCapability represents a skin type, e.g. StatTrak™ or Souvenir
type qualityCapability string

// skinnableItem represents any item that can be represented as a skin with a
// Market Hash Name. Internally, to derive the Market Hash Name, we require a
// descriptive name id for the language file, and whether the item is special.
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
	weapons map[string]*weapon
	gloves  map[string]*gloves
	crates  map[string]*itemCrate
}

// weapon represents a skinnable item that is also a weapon in csgo.
type weapon struct {
	id                    string
	languageNameId        string
	languageDescriptionId string
	special               bool
	prefab                *itemPrefab
}

func (w *weapon) getLanguageNameId() string {

	if w == nil {
		return ""
	}

	if w.languageNameId != "" {
		return w.languageNameId
	}

	return w.prefab.languageNameId
}

func (w *weapon) getLanguageDescriptionId() string {

	if w == nil {
		return ""
	}

	if w.languageDescriptionId != "" {
		return w.languageDescriptionId
	}

	return w.prefab.languageDescriptionId
}

func (w *weapon) getSpecial() bool {

	if w == nil {
		return false
	}

	return w.special
}

// mapToWeapon converts the provided map into a weapon providing
// all required parameters are present and of the correct type.
func mapToWeapon(data map[string]interface{}, prefabs map[string]*itemPrefab) (*weapon, error) {

	response := &weapon{}

	// get name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("id (name) missing from weapon")
	} else {
		response.id = val
	}

	// get language name id
	if val, err := crawlToType[string](data, "item_name"); err == nil {
		response.languageNameId = val
	}

	// get language description id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		response.languageDescriptionId = val
	}

	// get special
	if val, ok := data["prefab"].(string); ok {

		if val == "melee_unusual" {
			response.special = true
		}

		response.prefab = prefabs[val]
	}

	return response, nil
}

// gloves represents a special skinnable item that isn't a weapon.
type gloves struct {
	id                    string
	languageNameId        string
	languageDescriptionId string
}

func (g *gloves) getLanguageNameId() string {
	return g.languageNameId
}

func (g *gloves) getSpecial() bool {
	return true
}

// mapToGloves converts the provided map into gloves providing
// all required parameters are present and of the correct type.
func mapToGloves(data map[string]interface{}) (*gloves, error) {

	response := &gloves{}

	// get name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("id (name) missing from weapon") // TODO improve error
	} else {
		response.id = val
	}

	// get language name id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.New("language name id (item_name) missing from weapon") // TODO improve error
	} else {
		response.languageNameId = val
	}

	// get language description id
	if val, err := crawlToType[string](data, "item_description"); err != nil {
		return nil, errors.New("language description id (item_description) missing from weapon") // TODO improve error
	} else {
		response.languageDescriptionId = val
	}

	return response, nil
}

// itemCrate represents an openable crate that contains items. The crate's items
// are determined by the linked collection (item_set).
type itemCrate struct {
	id                    string
	languageNameId        string
	languageDescriptionId string

	// collectionId is the ID of the collection for the item/paintkit combinations
	// available in the crate.
	collectionId string

	// qualityCapability shows whether the crate can produce special skin qualities
	// e.g. Souvenir or StatTrak™
	qualityCapability qualityCapability
}

func (c *itemCrate) getLanguageNameId() string {
	return c.languageNameId
}

// mapToItemCrate converts the provided map into a itemCrate providing
// all required parameters are present and of the correct type.
func mapToItemCrate(data map[string]interface{}) (*itemCrate, error) {

	response := &itemCrate{
		qualityCapability: qualityNormal,
	}

	// get name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("id (name) missing from weapon") // TODO improve error
	} else {
		response.id = val
	}

	// get language name id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.New("language name id (item_name) missing from itemCrate") // TODO improve error
	} else {
		response.languageNameId = val
	}

	// get language description id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		response.languageDescriptionId = val
	}

	if val, err := crawlToType[string](data, "tags", "ItemSet", "tag_value"); err == nil {
		response.collectionId = val
	}

	return response, nil

}

// getItems processes the provided items data and, based on the item's prefab,
// produces the relevant item (e.g. gloves, weapon, or crate).
//
// All items are returned within the itemContainer part of the response.
func getItems(items map[string]interface{}, prefabs map[string]*itemPrefab) (*itemContainer, error) {

	response := &itemContainer{
		weapons: make(map[string]*weapon),
		gloves:  make(map[string]*gloves),
		crates:  make(map[string]*itemCrate),
	}

	items, err := crawlToType[map[string]interface{}](items, "items")
	if err != nil {
		return nil, errors.New("items missing from item data") // TODO format error better than this
	}

	for _, itemData := range items {

		itemMap, ok := itemData.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected item format found when fetching items")
		}

		prefab, ok := itemMap["prefab"].(string)
		if !ok {
			continue
		}

		switch getTypeFromPrefab(prefab, prefabs) {

		case itemTypeWeapon:
			w, err := mapToWeapon(itemMap, prefabs)
			if err != nil {
				return nil, err
			}

			response.weapons[w.id] = w

		case itemTypeGloves:
			g, err := mapToGloves(itemMap)
			if err != nil {
				return nil, err
			}

			response.gloves[g.id] = g

		case itemTypeCrate:
			c, err := mapToItemCrate(itemMap)
			if err != nil {
				return nil, err
			}

			response.crates[c.id] = c
		}
	}

	return response, nil
}
