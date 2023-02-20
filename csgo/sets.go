package csgo

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

// itemPaintkitSet is used to manage a weaponSet (set) of itemPaintKits such that
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

// weaponSet represents a weaponSet of items from the items_game file.
type weaponSet struct {
	id          string
	name        string
	description string
	items       map[string]string
}

// mapToWeaponSet converts the provided map into a weaponSet providing
// all required parameters are present and of the correct type.
//
// A response of nil, nil is returned when the provided set doesn't contain
// any weapons, e.g. a character set is provided.
func mapToWeaponSet(id string, data map[string]interface{}, language *language) (*weaponSet, error) {

	response := &weaponSet{
		id:    id,
		items: make(map[string]string),
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("language Name Id (name) missing from weaponSet")
	} else {

		lang, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "set_description"); err != nil {
		return nil, errors.New("language Description Id (description_string) missing from weaponSet")
	} else {

		lang, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.description = lang
	}

	items, err := crawlToType[map[string]interface{}](data, "items")
	if err != nil {
		return nil, errors.New("unable to find items in item_set")
	}

	for item, _ := range items {

		itemId, paintkitId, err := splitItemPaintkitString(item)
		if err != nil {
			continue
		}

		response.items[paintkitId] = itemId
	}

	// if set doesn't contain any weapons, return nothing
	if len(items) == 0 {
		return nil, nil
	}

	return response, nil
}

// splitItemPaintkitString splits a weaponSet item string that represents
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

// getWeaponSets will process all collections included in the provided items data
// (derived from items_game) and return them as a map[collectionId]*weaponSet.
func (c *csgoItems) getWeaponSets() (map[string]*weaponSet, error) {

	collections, err := crawlToType[map[string]interface{}](c.items, "item_sets")
	if err != nil {
		return nil, errors.New("item_sets missing from item data") // TODO handle error more nicely
	}

	response := make(map[string]*weaponSet)

	for setId, set := range collections {
		data, ok := set.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected format for item_set data")
		}

		setObj, err := mapToWeaponSet(setId, data, c.language)
		if err != nil {
			return nil, err
		}

		if setObj == nil {
			continue
		}

		response[setObj.id] = setObj
	}

	return response, nil
}

// characterSet represents a characterSet of items from the items_game file.
type characterSet struct {
	id          string
	name        string
	description string
	items       []string
}

// mapToCharacterSet converts the provided map into a characterSet providing
// all required parameters are present and of the correct type.
//
// A response of nil, nil is returned when the provided set doesn't contain
// any characters, e.g. a weapon set is provided.
func mapToCharacterSet(id string, data map[string]interface{}, language *language) (*characterSet, error) {

	response := &characterSet{
		id: id,
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("language Name Id (name) missing from weaponSet")
	} else {

		lang, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "set_description"); err != nil {
		return nil, errors.New("language Description Id (description_string) missing from weaponSet")
	} else {

		lang, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.description = lang
	}

	items, err := crawlToType[map[string]interface{}](data, "items")
	if err != nil {
		return nil, errors.New("unable to find items in item_set")
	}

	for item, _ := range items {

		if !strings.HasPrefix(item, "customplayer_") {
			continue
		}

		response.items = append(response.items, item)
	}

	// if set doesn't contain any characters, return nothing
	if len(items) == 0 {
		return nil, nil
	}

	return response, nil
}

// getCharacterSets will process all collections included in the provided items data
// (derived from items_game) and return them as a map[setId]*characterSet.
func (c *csgoItems) getCharacterSets() (map[string]*characterSet, error) {

	collections, err := crawlToType[map[string]interface{}](c.items, "item_sets")
	if err != nil {
		return nil, errors.New("item_sets missing from item data") // TODO handle error more nicely
	}

	response := make(map[string]*characterSet)

	for setId, set := range collections {
		data, ok := set.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected format for item_set data")
		}

		setObj, err := mapToCharacterSet(setId, data, c.language)
		if err != nil {
			return nil, err
		}

		if setObj == nil {
			continue
		}

		response[setObj.id] = setObj
	}

	return response, nil
}
