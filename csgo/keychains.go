package csgo

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type Keychain struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RarityId    string `json:"rarityId"`
}

func mapToKeychain(index int, data map[string]interface{}, language *language) (*Keychain, error) {
	response := &Keychain{
		Index: index,
	}

	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "Id (name) missing from Keychain")
	} else {
		response.Id = val
	}

	if val, err := crawlToType[string](data, "loc_name"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("loc_name missing from Keychain (%s)", response.Id))
	} else {
		// english file omits the # at the beggining for chains
		lang, _ := language.lookup(val[1:])
		response.Name = lang
	}

	if val, err := crawlToType[string](data, "loc_description"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("loc_description missing from Keychain (%s)", response.Id))
	} else {
		// english file omits the # at the beggining for chains
		lang, _ := language.lookup(val[1:])
		response.Description = lang
	}

	if val, err := crawlToType[string](data, "item_rarity"); err == nil {
		response.RarityId = val
	}

	return response, nil
}

func (c *csgoItems) getKeychains() (map[string]*Keychain, error) {
	response := make(map[string]*Keychain)

	chains, err := crawlToType[map[string]interface{}](c.items, "keychain_definitions")
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate keychain_definitions in provided items")
	}
	for index, chain := range chains {

		chainData, ok := chain.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Keychain data for %s is of unexpected type", index)
		}

		iIndex, err := strconv.Atoi(index)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret Keychain index (%s) as int", index))
		}

		chainMap, err := mapToKeychain(iIndex, chainData, c.language)

		response[chainMap.Id] = chainMap

	}

	return response, nil
}
