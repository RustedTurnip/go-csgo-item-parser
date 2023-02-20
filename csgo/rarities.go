package csgo

import "fmt"

// rarity represents a Csgo item rarity.
type rarity struct {
	id                  string
	generalRarityName   string
	weaponRarityName    string
	characterRarityName string
}

// mapToRarity converts the provided data map into a rarity object.
func mapToRarity(id string, data map[string]interface{}, language *language) (*rarity, error) {

	response := &rarity{
		id: id,
	}

	if key, err := crawlToType[string](data, "loc_key"); err == nil {
		name, err := language.lookup(key)
		if err == nil {
			response.generalRarityName = name
		}

	} else {
		return nil, fmt.Errorf("unable to locate language Name Id (loc_key) from Rarity %s: %s", response.id, err.Error())
	}

	if key, err := crawlToType[string](data, "loc_key_weapon"); err == nil {
		name, err := language.lookup(key)
		if err == nil {
			response.weaponRarityName = name
		}

	} else {
		return nil, fmt.Errorf("unable to locate language Name weapon Id (loc_key_weapon) from Rarity %s: %s", response.id, err.Error())
	}

	if key, err := crawlToType[string](data, "loc_key_character"); err == nil {
		name, err := language.lookup(key)
		if err == nil {
			response.weaponRarityName = name
		}

	} else {
		return nil, fmt.Errorf("unable to locate language Name character Id (loc_key_character) from Rarity %s: %s", response.id, err.Error())
	}

	return response, nil
}

// getRarities retrieves all Rarities from the provided items data and returns them
// in the format map[rarityId]Rarity.
func (c *csgoItems) getRarities() (map[string]*rarity, error) {

	response := make(map[string]*rarity)

	rarities, err := crawlToType[map[string]interface{}](c.items, "rarities")
	if err != nil {
		return nil, fmt.Errorf("unable to locate Rarities amongst items: %s", err.Error())
	}

	for id, rarity := range rarities {

		rarityData, ok := rarity.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Rarity data for %s is of unexpected type", id)
		}

		rarityMap, err := mapToRarity(id, rarityData, c.language)
		if err != nil {
			return nil, err
		}

		response[id] = rarityMap
	}

	return response, nil
}
