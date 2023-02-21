package csgo

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

var (
	// weaponPaintkitRe is the pattern of a item Id and paintkit Id item set
	// string that looks like: "[paint_kit_id]weapon_id"
	weaponPaintkitRe = regexp.MustCompile("^\\[([a-zA-Z0-9_\\-)]+)\\]([a-zA-Z0-9_\\-]+)$")
)

// WeaponSet represents a WeaponSet of items from the items_game file.
type WeaponSet struct {
	Id          string
	Name        string
	Description string
	Items       map[string]string
}

// mapToWeaponSet converts the provided map into a WeaponSet providing
// all required parameters are present and of the correct type.
//
// A response of nil, nil is returned when the provided set doesn't contain
// any weapons, e.g. a character set is provided.
func mapToWeaponSet(id string, data map[string]interface{}, language *language) (*WeaponSet, error) {

	response := &WeaponSet{
		Id:    id,
		Items: make(map[string]string),
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "language Name Id (name) missing from WeaponSet")
	} else {

		lang, err := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to lookup WeaponSet's name (%s)", val))
		}

		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "set_description"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("language Description Id (description_string) missing from WeaponSet (%s)", response.Id))
	} else {

		lang, err := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to lookup WeaponSet's description (%s)", val))
		}

		response.Description = lang
	}

	items, err := crawlToType[map[string]interface{}](data, "items")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to find items in item_set %s", response.Id))
	}

	for item, _ := range items {

		itemId, paintkitId, err := splitItemPaintkitString(item)
		if err != nil {
			continue
		}

		response.Items[paintkitId] = itemId
	}

	// if set doesn't contain any weapons, return nothing
	if len(items) == 0 {
		return nil, nil
	}

	return response, nil
}

// splitItemPaintkitString splits a WeaponSet item string that represents
// an item ID - paintkit ID mapping into an item ID and paintkit ID.
//
// If the provided string cannot be parsed into the two ids, an error is
// returned.
func splitItemPaintkitString(itemPaintkit string) (string, string, error) {

	match := weaponPaintkitRe.FindStringSubmatch(itemPaintkit)
	if len(match) != 3 {
		fmt.Println(itemPaintkit)
		return "", "", errors.New("unexpected [Weapon]paintkit format")
	}

	return match[2], match[1], nil
}

// getWeaponSets will process all collections included in the provided items data
// (derived from items_game) and return them as a map[collectionId]*WeaponSet.
func (c *csgoItems) getWeaponSets() (map[string]*WeaponSet, error) {

	collections, err := crawlToType[map[string]interface{}](c.items, "item_sets")
	if err != nil {
		return nil, errors.Wrap(err, "item_sets missing from item data")
	}

	response := make(map[string]*WeaponSet)

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

		response[setObj.Id] = setObj
	}

	return response, nil
}

// CharacterSet represents a CharacterSet of items from the items_game file.
type CharacterSet struct {
	Id          string
	Name        string
	Description string
	Items       []string
}

// mapToCharacterSet converts the provided map into a CharacterSet providing
// all required parameters are present and of the correct type.
//
// A response of nil, nil is returned when the provided set doesn't contain
// any characters, e.g. a Weapon set is provided.
func mapToCharacterSet(id string, data map[string]interface{}, language *language) (*CharacterSet, error) {

	response := &CharacterSet{
		Id: id,
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "language Name Id (name) missing from CharacterSet")
	} else {

		lang, _ := language.lookup(val)
		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "set_description"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("language Description Id (description_string) missing from CharacterSet (%s)", response.Id))
	} else {

		lang, _ := language.lookup(val)
		response.Description = lang
	}

	items, err := crawlToType[map[string]interface{}](data, "items")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to find items in CharacterSet (%s)", response.Id))
	}

	for item, _ := range items {

		if !strings.HasPrefix(item, "customplayer_") {
			continue
		}

		response.Items = append(response.Items, item)
	}

	// if set doesn't contain any characters, return nothing
	if len(items) == 0 {
		return nil, nil
	}

	return response, nil
}

// getCharacterSets will process all collections included in the provided items data
// (derived from items_game) and return them as a map[setId]*CharacterSet.
func (c *csgoItems) getCharacterSets() (map[string]*CharacterSet, error) {

	collections, err := crawlToType[map[string]interface{}](c.items, "item_sets")
	if err != nil {
		return nil, errors.Wrap(err, "item_sets missing from item data")
	}

	response := make(map[string]*CharacterSet)

	for setId, set := range collections {
		data, ok := set.(map[string]interface{})
		if !ok {
			return nil, errors.Wrap(err, "unexpected format for item_set data")
		}

		setObj, err := mapToCharacterSet(setId, data, c.language)
		if err != nil {
			return nil, err
		}

		if setObj == nil {
			continue
		}

		response[setObj.Id] = setObj
	}

	return response, nil
}