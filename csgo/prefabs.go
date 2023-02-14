package csgo

import "errors"

var (
	// itemPrefabPrefabs is a map of all prefabs that exist against item prefabs
	// that we need to track.
	itemPrefabPrefabs = map[string]itemType{

		"": itemTypeUnknown,

		// weapons
		"primary":       itemTypeWeapon, // covers ay weapon that can be primary (e.g. smg)
		"secondary":     itemTypeWeapon, // covers any weapon that can be secondary (e.g. pistol)
		"melee_unusual": itemTypeWeapon, // covers all tradable knives

		// gloves
		"hands": itemTypeGloves, // covers gloves

		// crates
		"weapon_case":             itemTypeCrate,
		"weapon_case_souvenirpkg": itemTypeCrate,
	}
)

// itemPrefab represents a csgo prefab which is used to categorise item
// types. e.g. melee, primary (both of which are weapons).
type itemPrefab struct {
	id                    string
	parentPrefab          string
	languageNameId        string
	languageDescriptionId string
	special               bool
	itemType              itemType
}

// mapToItemPrefab converts the provided map (data) into a prefab object.
func mapToItemPrefab(id string, data map[string]interface{}) *itemPrefab {

	response := &itemPrefab{
		id:       id,
		itemType: itemPrefabPrefabs[id],
	}

	if val, ok := data["prefab"].(string); ok {
		response.parentPrefab = val
	}

	if val, ok := data["item_name"].(string); ok {
		response.languageNameId = val
	}

	if val, ok := data["item_description"].(string); ok {
		response.languageDescriptionId = val
	}

	if val, ok := data["prefab"].(string); ok {
		if val == "melee_unusual" {
			response.special = true
		}
	}

	return response
}

// getItemPrefabs retrieves all required prefabs from the provided items
// map and returns them in the format map[prefabId]itemPrefab.
func getItemPrefabs(items map[string]interface{}) (map[string]*itemPrefab, error) {

	response := make(map[string]*itemPrefab)

	prefabs, err := crawlToType[map[string]interface{}](items, "prefabs")
	if err != nil {
		return nil, errors.New("item data is missing prefabs") // TODO improve error
	}

	for prefabId, prefab := range prefabs {

		prefabData, ok := prefab.(map[string]interface{})
		if !ok {
			return nil, errors.New("prefab in in unexpected format")
		}

		// build prefab
		obj := mapToItemPrefab(prefabId, prefabData)
		response[obj.id] = obj
	}

	return response, nil
}

// getTypeFromPrefab will traverse the prefab inheritence tree until a
// recognisable prefab is found (or the end is reached) and it will
// return the associated item type for the matched prefab.
func getTypeFromPrefab(prefabId string, prefabs map[string]*itemPrefab) itemType {

	// if provided prefab id is of recognised type, return type
	if prefabType, ok := itemPrefabPrefabs[prefabId]; ok {
		return prefabType
	}

	// if no parent prefab to lookup, return unknown
	if prefab, ok := prefabs[prefabId]; ok {
		// crawl further
		return getTypeFromPrefab(prefab.parentPrefab, prefabs)
	}

	// else type is unknown so return that
	return itemTypeUnknown
}
