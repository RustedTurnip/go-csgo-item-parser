package csgo

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	weaponPaintkitRe = regexp.MustCompile("^\\[([a-zA-Z0-9_\\-)]+)\\]([a-zA-Z0-9_\\-]+)$")
)

type itemPaintkit struct {
	itemId     string
	paintkitId string
}

// TODO comment
type itemPaintkitSet struct {
	itemPaintkits map[string]*itemPaintkit
}

// TODO comment
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

// TODO comment
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

// mapToWeaponSet converts the provided map into a collection providing
// all required parameters are present and of the correct type.
func mapToWeaponSet(id string, data map[string]interface{}) (*collection, error) {

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

// TODO comment
func splitItemPaintkitString(itemPaintkit string) (string, string, error) {

	match := weaponPaintkitRe.FindStringSubmatch(itemPaintkit)
	if len(match) != 3 {
		fmt.Println(itemPaintkit)
		return "", "", errors.New("unexpected [weapon]paintkit format")
	}

	return match[2], match[1], nil
}

func getWeaponSets(items map[string]interface{}) (map[string]*collection, error) {

	sets, err := crawlToType[map[string]interface{}](items, "item_sets")
	if err != nil {
		return nil, errors.New("item_sets missing from item data") // TODO handle error more nicely
	}

	response := make(map[string]*collection)

	for setId, set := range sets {
		data, ok := set.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected format for item_set data")
		}

		setObj, err := mapToWeaponSet(setId, data)
		if err != nil {
			return nil, err
		}

		response[setObj.id] = setObj
	}

	return response, nil
}
