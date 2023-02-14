package csgo

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	// weaponPaintkitRe is the pattern of a item Id and paintkit Id item set
	// string that looks like: "[paint_kit_id]weapon_id"
	weaponPaintkitRe = regexp.MustCompile("^\\[([a-zA-Z0-9_\\-)]+)\\]([a-zA-Z0-9_\\-]+)$")
)

// itemPaintkit represents an itemId paintkitId pairing.
type itemPaintkit struct {
	itemId     string
	paintkitId string
}

// itemPaintkitSet is used to manage a collection (set) of itemPaintKits such that
// any pair is only stored once.
type itemPaintkitSet struct {
	itemPaintkits map[string]*itemPaintkit
}

// add stores an itemPaintkit pair in the itemPaintkitSet if the same pair
// hasn't already been added (otherwise, noop).
func (ip *itemPaintkitSet) add(itemId, paintkitId string) {

	// serialise ids to create unique identifier
	id := fmt.Sprintf("%s_%s", itemId, paintkitId)

	if ip.itemPaintkits == nil {
		ip.itemPaintkits = make(map[string]*itemPaintkit)
	}

	// if combination already, no need to add
	if _, ok := ip.itemPaintkits[id]; ok {
		return
	}

	// store combination
	ip.itemPaintkits[id] = &itemPaintkit{
		itemId:     itemId,
		paintkitId: paintkitId,
	}
}

// forEachItemPaintkit provides a way to iterate over the contents of a
// itemPaintkitSet where fn allows you to provide the action to be performed
// upon each itemPaintkit pair.
//
// If fn returns an error, the loop will break early and the error is bubbled
// up and returned by forEachItemPaintkit.
func (ip *itemPaintkitSet) forEachItemPaintkit(fn func(itemId, paintkitId string) error) error {

	for _, itemPaintkit := range ip.itemPaintkits {
		err := fn(itemPaintkit.itemId, itemPaintkit.paintkitId)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO comment struct
type collection struct {
	id                    string
	languageNameId        string
	languageDescriptionId string
	items                 *itemPaintkitSet
}

// mapToCollection converts the provided map into a collection providing
// all required parameters are present and of the correct type.
func mapToCollection(id string, data map[string]interface{}) (*collection, error) {

	response := &collection{
		id:    id,
		items: &itemPaintkitSet{},
	}

	// get language name id
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("language name id (name) missing from collection")
	} else {
		response.languageNameId = val
	}

	// get language description id
	if val, err := crawlToType[string](data, "set_description"); err != nil {
		return nil, errors.New("language description id (description_string) missing from collection")
	} else {
		response.languageDescriptionId = val
	}

	items, err := crawlToType[map[string]interface{}](data, "items")
	if err != nil {
		return nil, errors.New("unable to find items in item_set")
	}

	for item, _ := range items {

		itemId, paintkitId, err := splitItemPaintkitString(item)
		if err != nil {
			// TODO this will ignore any set items that are not weapons
			continue
		}

		response.items.add(itemId, paintkitId)
	}

	return response, nil
}

// splitItemPaintkitString splits a collection item string that represents
// an item ID - paintkit ID mapping into an item ID and paintkit ID.
//
// If the provided string cannot be parsed into the two ids, an error is
// returned.
func splitItemPaintkitString(itemPaintkit string) (string, string, error) {

	match := weaponPaintkitRe.FindStringSubmatch(itemPaintkit)
	if len(match) != 3 {
		fmt.Println(itemPaintkit)
		return "", "", errors.New("unexpected [weapon]paintkit format")
	}

	return match[2], match[1], nil
}

// getCollections will process all collections included in the provided items data
// (derived from items_game) and return them as a map[collectionId]*collection.
func getCollections(items map[string]interface{}) (map[string]*collection, error) {

	collections, err := crawlToType[map[string]interface{}](items, "item_sets")
	if err != nil {
		return nil, errors.New("item_sets missing from item data") // TODO handle error more nicely
	}

	response := make(map[string]*collection)

	for setId, set := range collections {
		data, ok := set.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected format for item_set data")
		}

		setObj, err := mapToCollection(setId, data)
		if err != nil {
			return nil, err
		}

		response[setObj.id] = setObj
	}

	return response, nil
}
