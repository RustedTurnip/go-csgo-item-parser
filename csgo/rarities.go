package csgo

import "fmt"

// TODO comment this
type rarity struct {
	id                      string
	languageNameId          string
	languageNameWeaponId    string
	languageNameCharacterId string
	colourId                string
}

// TODO comment
func mapToRarity(id string, data map[string]interface{}) (*rarity, error) {

	response := &rarity{
		id: id,
	}

	if key, err := crawlToType[string](data, "loc_key"); err == nil {
		response.languageNameId = key
	} else {
		return nil, fmt.Errorf("unable to locate language name id (loc_key) from rarity %s: %s", response.id, err.Error())
	}

	if key, err := crawlToType[string](data, "loc_key_weapon"); err == nil {
		response.languageNameWeaponId = key
	} else {
		return nil, fmt.Errorf("unable to locate language name weapon id (loc_key_weapon) from rarity %s: %s", response.id, err.Error())
	}

	if key, err := crawlToType[string](data, "loc_key_character"); err == nil {
		response.languageNameWeaponId = key
	} else {
		return nil, fmt.Errorf("unable to locate language name character id (loc_key_character) from rarity %s: %s", response.id, err.Error())
	}

	if key, err := crawlToType[string](data, "color"); err == nil {
		response.colourId = key
	} else {
		return nil, fmt.Errorf("unable to locate colour id (color) from rarity %s: %s", response.id, err.Error())
	}

	return response, nil
}

// TODO comment
func getRarities(items map[string]interface{}) (map[string]*rarity, error) {

	response := make(map[string]*rarity)

	rarities, err := crawlToType[map[string]interface{}](items, "rarities")
	if err != nil {
		return nil, fmt.Errorf("unable to locate rarities amongst items: %s", err.Error())
	}

	for id, rarity := range rarities {

		rarityData, ok := rarity.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("rarity data for %s is of unexpected type", id)
		}

		rarityMap, err := mapToRarity(id, rarityData)
		if err != nil {
			return nil, err
		}

		response[id] = rarityMap
	}

	return response, nil
}
