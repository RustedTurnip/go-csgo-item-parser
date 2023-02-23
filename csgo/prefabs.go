package csgo

import (
	"github.com/pkg/errors"
)

// itemPrefab represents a Csgo prefab which is used to categorise item
// types. e.g. melee, primary (both of which are Guns).
type itemPrefab struct {
	id           string
	parentPrefab string
	name         string
	description  string
}

// mapToItemPrefab converts the provided map (data) into a prefab object.
func mapToItemPrefab(id string, data map[string]interface{}, language *language) (*itemPrefab, error) {

	response := &itemPrefab{
		id: id,
	}

	if val, ok := data["prefab"].(string); ok {
		response.parentPrefab = val
	}

	if val, ok := data["item_name"].(string); ok {
		lang, _ := language.lookup(val)
		response.name = lang
	}

	if val, ok := data["item_description"].(string); ok {
		lang, _ := language.lookup(val)
		response.description = lang
	}

	return response, nil
}

// getItemPrefabs retrieves all required prefabs from the provided items
// map and returns them in the format map[prefabId]itemPrefab.
func (c *csgoItems) getItemPrefabs() (map[string]*itemPrefab, error) {

	response := make(map[string]*itemPrefab)

	prefabs, err := crawlToType[map[string]interface{}](c.items, "prefabs")
	if err != nil {
		return nil, errors.Wrap(err, "item data is missing prefabs")
	}

	for prefabId, prefab := range prefabs {

		prefabData, ok := prefab.(map[string]interface{})
		if !ok {
			return nil, errors.New("prefab in in unexpected format")
		}

		// build prefab
		obj, err := mapToItemPrefab(prefabId, prefabData, c.language)
		if err != nil {
			return nil, err
		}

		response[obj.id] = obj
	}

	return response, nil
}
