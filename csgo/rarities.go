package csgo

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

// Rarity represents a Csgo item Rarity.
type Rarity struct {
	Id                  string
	Index               int
	GeneralRarityName   string
	WeaponRarityName    string
	CharacterRarityName string
}

// mapToRarity converts the provided data map into a Rarity object.
func mapToRarity(id string, data map[string]interface{}, language *language) (*Rarity, error) {

	response := &Rarity{
		Id: id,
	}

	// get index
	if val, ok := data["value"].(string); ok {
		if valInt, err := strconv.Atoi(val); err == nil {
			response.Index = valInt
		} else {
			return nil, errors.Wrap(err, fmt.Sprintf("unexpected index (value) type: %s", val))
		}
	} else {
		return nil, fmt.Errorf("rarity (%s) missing expected field \"value\"", response.Id)
	}

	if key, err := crawlToType[string](data, "loc_key"); err == nil {
		lang, _ := language.lookup(key)
		response.GeneralRarityName = lang
	} else {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to locate language Name Id (loc_key) from Rarity %s", response.Id))
	}

	if key, err := crawlToType[string](data, "loc_key_weapon"); err == nil {
		lang, _ := language.lookup(key)
		response.WeaponRarityName = lang
	} else {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to locate language Name Weapon Id (loc_key_weapon) from Rarity %s", response.Id))
	}

	if key, err := crawlToType[string](data, "loc_key_character"); err == nil {
		lang, _ := language.lookup(key)
		response.WeaponRarityName = lang
	} else {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to locate language Name character Id (loc_key_character) from Rarity %s", response.Id))
	}

	return response, nil
}

// getRarities retrieves all Rarities from the provided items data and returns them
// in the format map[rarityId]Rarity.
func (c *csgoItems) getRarities() (map[string]*Rarity, error) {

	response := make(map[string]*Rarity)

	rarities, err := crawlToType[map[string]interface{}](c.items, "rarities")
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate rarities amongst items")
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
