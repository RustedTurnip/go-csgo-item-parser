package csgo

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// itemPrefabPrefabs is a map of all prefabs that exist against item prefabs
	// that we need to track.
	itemPrefabPrefabs = map[string]itemType{

		"": itemTypeUnknown,

		// Guns
		"primary":       itemTypeWeaponGun,   // covers ay weapon that can be primary (e.g. smg)
		"secondary":     itemTypeWeaponGun,   // covers any weapon that can be secondary (e.g. pistol)
		"melee_unusual": itemTypeWeaponKnife, // covers all tradable Knives

		// Gloves
		"hands": itemTypeGloves, // covers Gloves

		// crates
		"weapon_case":             itemTypeCrate,
		"weapon_case_souvenirpkg": itemTypeCrate,

		// stickers
		"sticker_capsule": itemTypeStickerCapsule,
	}
)

// itemPrefab represents a Csgo prefab which is used to categorise item
// types. e.g. melee, primary (both of which are Guns).
type itemPrefab struct {
	id           string
	parentPrefab string
	name         string
	description  string
	itemType     itemType
}

// mapToItemPrefab converts the provided map (data) into a prefab object.
func mapToItemPrefab(id string, data map[string]interface{}, language *language) (*itemPrefab, error) {

	if id == "atlanta2017_sticker_capsule_prefab" {
		fmt.Println()
	}

	response := &itemPrefab{
		id:       id,
		itemType: itemPrefabPrefabs[id],
	}

	if val, ok := data["prefab"].(string); ok {
		response.parentPrefab = val
	}

	if val, ok := data["item_name"].(string); ok {
		lang, err := language.lookup(val)
		if err == nil {
			response.name = lang
		}
	}

	if val, ok := data["item_description"].(string); ok {
		lang, err := language.lookup(val)
		if err == nil {
			response.description = lang
		}
	}

	response.itemType = getPrefabItemType(id, data)

	return response, nil
}

// getItemPrefabs retrieves all required prefabs from the provided items
// map and returns them in the format map[prefabId]itemPrefab.
func (c *csgoItems) getItemPrefabs() (map[string]*itemPrefab, error) {

	response := make(map[string]*itemPrefab)

	prefabs, err := crawlToType[map[string]interface{}](c.items, "prefabs")
	if err != nil {
		return nil, errors.New("item data is missing prefabs") // TODO improve error
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

// getTypeFromPrefab will traverse the prefab inheritence tree until a
// recognisable prefab is found (or the end is reached) and it will
// return the associated item type for the matched prefab.
func getTypeFromPrefab(prefabId string, prefabs map[string]*itemPrefab) itemType {

	prefabId = strings.Split(prefabId, " ")[0]

	prefab, ok := prefabs[prefabId]

	// if prefab isn't recognised
	if !ok {
		return itemTypeUnknown
	}

	// if prefab item type is known, return item type
	if prefab.itemType != itemTypeUnknown {
		return prefab.itemType
	}

	// if at root of prefab tree, return unknown
	if prefab.parentPrefab == "" {
		return itemTypeUnknown
	}

	// if prefab item type is unknown, but prefab has parent, crawl further
	return getTypeFromPrefab(prefab.parentPrefab, prefabs)
}

// getPrefabItemType can be used to identify some prefabs for the itemType they
// represent. There isn't a consistent way to do this, but a combination of checks
// are used to identify the item type.
func getPrefabItemType(id string, data map[string]interface{}) itemType {

	if it, ok := itemPrefabPrefabs[id]; ok {
		return it
	}

	if val, _ := data["inv_container_and_tools"].(string); val == "sticker_capsule" {
		return itemTypeStickerCapsule
	}

	if _, err := crawlToType[map[string]interface{}](data, "tags", "StickerCapsule"); err == nil {
		return itemTypeStickerCapsule
	}

	return itemTypeUnknown
}
