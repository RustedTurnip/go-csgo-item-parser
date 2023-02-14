package csgo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rustedturnip/go-csgo-item-parser/entities"
)

// GetAllItems is the main entrypoint of the csgo package and is responsible for
// parsing the language and item data, and transforming it into the universal
// entities.
func GetAllItems(languageData, itemData map[string]interface{}) (*entities.Items, error) {

	csgo, err := newCsgo(languageData, itemData)
	if err != nil {
		return nil, err
	}

	return csgo.getAllItems()
}

// language represents a csgo language file that provides the descriptions
// and descriptive names of in game items.
type language struct {
	data map[string]interface{}
}

// lookup will return the string value (e.g. descriptive name) of the
// provided identifier (key).
//
// If the key cannot be found, an error is returned.
func (l *language) lookup(key string) (string, error) {

	// remove pound sign from beginning of string
	key = strings.TrimPrefix(key, "#")

	val, err := crawlToType[string](l.data, strings.ToLower(key))
	if err != nil {
		return "", fmt.Errorf("could not locate token (%s) in language: %s", key, err.Error())
	}

	return val, nil
}

// newLanguage takes in the data map of an already parsed language file and
// returns a language "client" that can be used to perform key lookups.
func newLanguage(data map[string]interface{}) (*language, error) {

	// check language base data exists
	lang, ok := data["lang"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to locate \"lang\" in provided languageData")
	}

	// check lang has tokens as expected
	tokens, ok := lang["Tokens"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to locate \"lang/Tokens\" in provided languageData")
	}

	l := &language{
		data: make(map[string]interface{}),
	}

	for k, v := range tokens {
		l.data[strings.ToLower(k)] = v
	}

	return l, nil
}

// csgo is a representation of all csgo items that are relevant to interpreting
// the game_items file.
type csgo struct {
	language *language

	// custom types
	prefabs     map[string]*itemPrefab
	rarities    map[string]*rarity
	paintkits   map[string]*paintkit
	stickerkits map[string]*stickerkit
	weaponSets  map[string]*collection
	knifeSet    *itemPaintkitSet
	gloveSet    *itemPaintkitSet

	// items
	weapons      map[string]*weapon
	gloves       map[string]*gloves
	weaponCrates map[string]*itemCrate
}

// getAllItems will perform the necessary joins of the already parsed
// items from the items_game file and return more universally recognised
// entities like skin, stickers etc.
func (c *csgo) getAllItems() (*entities.Items, error) {
	items := &entities.Items{}

	skins, err := c.getAllSkins()
	if err != nil {
		return nil, err
	}

	items.Skins = skins

	return items, nil
}

// newCsgo represents the constructor for csgo and will perform the necessary
// preprocessing of the language and item data.
func newCsgo(languageData, itemData map[string]interface{}) (*csgo, error) {

	language, err := newLanguage(languageData)
	if err != nil {
		return nil, err
	}

	// check items base data exists
	fileItems, err := crawlToType[map[string]interface{}](itemData, "items_game")
	if err != nil {
		return nil, errors.New("unable to locate \"items_game\" in provided itemData") // TODO format error better than this
	}

	// build entities from file
	prefabs, err := getItemPrefabs(fileItems)
	if err != nil {
		return nil, err
	}

	rarities, err := getRarities(fileItems)
	if err != nil {
		return nil, err
	}

	paintkits, err := getPaintkits(fileItems)
	if err != nil {
		return nil, err
	}

	stickerkits, err := getStickerkits(fileItems)
	if err != nil {
		return nil, err
	}

	sets, err := getCollections(fileItems)
	if err != nil {
		return nil, err
	}

	items, err := getItems(fileItems, prefabs)
	if err != nil {
		return nil, err
	}

	// knives are not categorised into sets within the items_game.txt file,
	// so they are handled separately.
	knifeSet, err := getIconSet(fileItems, items.weapons)
	if err != nil {
		return nil, err
	}

	// gloves are not categorised into sets within the items_game.txt file,
	// so they are handled separately.
	gloveSet, err := getIconSet(fileItems, items.gloves)
	if err != nil {
		return nil, err
	}

	return &csgo{
		language: language,

		prefabs:     prefabs,
		rarities:    rarities,
		paintkits:   paintkits,
		stickerkits: stickerkits,
		weaponSets:  sets,
		knifeSet:    knifeSet,
		gloveSet:    gloveSet,

		weapons:      items.weapons,
		gloves:       items.gloves,
		weaponCrates: items.crates,
	}, nil
}

// crawlToType will traverse down the provided map (m) through the provided
// path of keys (path) until reaching the last node in the path. It returns
// the last node as the provided type (T), however will return an error if
// the node can't be cast to that type.
//
// Note: it is assumed that all but the final node in the path is of type
// map[string]interface{}
func crawlToType[T any](m map[string]interface{}, path ...string) (T, error) {

	var empty T

	if len(path) == 1 {
		return crawl[T](m, path[0])
	}

	nested, err := crawl[map[string]interface{}](m, path[0])
	if err != nil {
		return empty, err
	}

	return crawlToType[T](nested, path[1:]...)
	// TODO return specific errors i.e. errNotFound and errUnexpectedType
}

// crawl will return the value at the provided key casting it to the provided
// type (T). If the value doesn't exist or doesn't match the type provided, an
// error will be returned.
func crawl[T any](m map[string]interface{}, key string) (T, error) {

	var empty T // equivalent of nil

	val, ok := m[key].(T)
	if !ok {
		return empty, fmt.Errorf("couldn't find key %s in provided map", key)
	}

	return val, nil
}
