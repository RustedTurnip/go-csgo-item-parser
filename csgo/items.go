package csgo

import (
	"errors"
)

type qualityCapability string

type skinableItem interface {
	getLanguageNameId() string
	getSpecial() bool
}

var (
	qualityNormal   qualityCapability = "Normal"
	qualityStatTrak qualityCapability = "StatTrak™"
	qualitySouvenir qualityCapability = "Souvenir"
)

// TODO comment struct
// Gloves detect because "item_name" starts with "#CSGO_Wearable_"
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

// TODO comment here
type weaponCrate struct {
	id                    string
	languageNameId        string
	languageDescriptionId string
	weaponSetId           string

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

	// get name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("id (name) missing from weapon") // TODO improve error
	} else {
		response.id = val
	}

	// get language name id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.New("language name id (item_name) missing from weaponCrate") // TODO improve error
	} else {
		response.languageNameId = val
	}

	// get language description id
	if val, err := crawlToType[string](data, "item_description"); err == nil {
		response.languageDescriptionId = val
	}

	if val, err := crawlToType[string](data, "tags", "ItemSet", "tag_value"); err == nil {
		response.weaponSetId = val
	}

	return response, nil

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

type itemContainer struct {
	weapons map[string]*weapon
	gloves  map[string]*gloves
	crates  map[string]*weaponCrate
}

func getItems(items map[string]interface{}, prefabs map[string]*itemPrefab) (*itemContainer, error) {

	response := &itemContainer{
		weapons: make(map[string]*weapon),
		gloves:  make(map[string]*gloves),
		crates:  make(map[string]*weaponCrate),
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
			c, err := mapToWeaponCrate(itemMap)
			if err != nil {
				return nil, err
			}

			response.crates[c.id] = c
		}
	}

	return response, nil
}
