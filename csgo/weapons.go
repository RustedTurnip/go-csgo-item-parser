package csgo

import (
	"errors"
)

// TODO comment struct
// Weapon detect because "item_name" starts with "#SFUI_WPNHUD_" (includes knifes and guns)
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
