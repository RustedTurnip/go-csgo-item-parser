package csgo

import "errors"

type itemType int

const (
	itemTypeUnknown itemType = iota
	itemTypeWeapon
	itemTypeGloves
	itemTypeCrate
)

var (
	// weaponPrefabPrefabs is a map of all prefabs that exist against weapon prefabs
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

type itemPrefab struct {
	id                    string
	parentPrefab          string
	languageNameId        string
	languageDescriptionId string
	special               bool
	itemType              itemType
}

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
