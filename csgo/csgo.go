package csgo

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// New takes the required languageData and itemData maps (from csgo_english.txt and
// items_game.txt respectively) and extracts the desired sub elements from them,
// returning a fully instantiated Csgo.
func New(languageData, itemData map[string]interface{}) (*Csgo, error) {
	language, err := newLanguage(languageData)
	if err != nil {
		return nil, err
	}

	// check items base data exists
	fileItems, err := crawlToType[map[string]interface{}](itemData, "items_game")
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate \"items_game\" in provided itemData")
	}

	items, err := newCsgoItems(fileItems, language)
	if err != nil {
		return nil, err
	}

	rarities, err := items.getRarities()
	if err != nil {
		return nil, err
	}

	qualities, err := items.getQualities()
	if err != nil {
		return nil, err
	}

	paintkits, err := items.getPaintkits()
	if err != nil {
		return nil, err
	}

	keychains, err := items.getKeychains()
	if err != nil {
		return nil, err
	}

	musickits, err := items.getMusickits()
	if err != nil {
		return nil, err
	}

	stickerEnteties, err := items.getStickerkits()
	if err != nil {
		return nil, err
	}

	weaponSets, err := items.getWeaponSets()
	if err != nil {
		return nil, err
	}

	itemEntities, err := items.getItems()
	if err != nil {
		return nil, err
	}

	// Knives are not categorised into sets within the items_game.txt file,
	// so they are handled separately.
	knifeSet, err := items.getKnifeSet(mapTypeToMapInterface(itemEntities.knives))
	if err != nil {
		return nil, err
	}

	// Gloves are not categorised into sets within the items_game.txt file,
	// so they are handled separately.
	gloveSet, err := items.getIconSet(mapTypeToMapInterface(itemEntities.gloves))
	if err != nil {
		return nil, err
	}

	return &Csgo{
		Rarities:   rarities,
		Qualities:  qualities,
		Paintkits:  paintkits,
		Keychains:  keychains,
		Musickits:  musickits,
		WeaponSets: weaponSets,
		KnifeSet:   knifeSet,
		GloveSet:   gloveSet,

		Stickerkits: stickerEnteties.stickers,
		Spraykits:   stickerEnteties.sprays,
		Patchkits:   stickerEnteties.patches,

		Guns:            itemEntities.weapons,
		Knives:          itemEntities.knives,
		Gloves:          itemEntities.gloves,
		Equipment:       itemEntities.equipment,
		WeaponCrates:    itemEntities.crates,
		StickerCapsules: itemEntities.stickerCapsules,
		Tools:           itemEntities.tools,
		Characters:      itemEntities.characters,
		Collectables:    itemEntities.collectables,
	}, nil
}

// language represents a Csgo language file that provides the descriptions
// and descriptive names of in game items.
type language struct {
	data map[string]interface{}
}

// lookup will return the string value (e.g. descriptive Name) of the
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

// csgoItems provides a wrapper for the data from csgo_english.txt and items_game.txt.
type csgoItems struct {
	items    map[string]interface{}
	language *language

	// cache attributes
	prefabs            map[string]*itemPrefab
	revolvingLootLists revolvingLootLists
	clientLootLists    map[string]*clientLootList
}

// newCsgoItems is the csgoItems constructor.
func newCsgoItems(itemData map[string]interface{}, language *language) (*csgoItems, error) {

	response := &csgoItems{
		items:    itemData,
		language: language,
	}

	prefabs, err := response.getItemPrefabs()
	if err != nil {
		return nil, err
	}

	response.prefabs = prefabs

	revolvingLootLists, err := response.getRevolvingLootLists()
	if err != nil {
		return nil, err
	}

	response.revolvingLootLists = revolvingLootLists

	clientLootLists, err := response.getClientLootLists()
	if err != nil {
		return nil, err
	}

	response.clientLootLists = clientLootLists

	return response, nil
}

// Csgo is a representation of all Csgo items that are relevant to interpreting
// the game_items file.
type Csgo struct {

	// CSGO types
	Rarities   map[string]*Rarity    `json:"Rarities"`
	Qualities  map[string]*Quality   `json:"Qualities"`
	Paintkits  map[string]*Paintkit  `json:"Paintkits"`
	Keychains  map[string]*Keychain  `json:"Keychains"`
	Musickits  map[string]*Musickit  `json:"Musickit"`
	WeaponSets map[string]*WeaponSet `json:"WeaponSets"`
	KnifeSet   map[string][]string   `json:"KnifeSet"`
	GloveSet   map[string][]string   `json:"GloveSet"`

	// Sticker subtypes
	Stickerkits map[string]*Stickerkit `json:"Stickerkits"`
	Spraykits   map[string]*Spraykit   `json:"Spraykits"`
	Patchkits   map[string]*Patchkit   `json:"Patchkits"`

	// items
	Guns            map[string]*Weapon         `json:"Guns"`
	Knives          map[string]*Weapon         `json:"Knives"`
	Gloves          map[string]*Gloves         `json:"Gloves"`
	Equipment       map[string]*Equipment      `json:"Equipment"`
	Tools           map[string]*Tool           `json:"Tools"`
	WeaponCrates    map[string]*WeaponCrate    `json:"WeaponCrates"`
	StickerCapsules map[string]*StickerCapsule `json:"StickerCapsules"`
	Characters      map[string]*Character      `json:"Characters"`
	// some might not have descriptions due to them being placeholders
	Collectables map[string]*Collectible `json:"Collectables"`
}

var (
	// errCrawlNotFound is used to distinguish an error from not being able to
	// find the node at the provided path.
	errCrawlNotFound = errors.New("the node at the provided path could not be found")
)

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
}

// crawl will return the value at the provided key casting it to the provided
// type (T). If the value doesn't exist or doesn't match the type provided, an
// error will be returned.
func crawl[T any](m map[string]interface{}, key string) (T, error) {

	var empty T // equivalent of nil

	val, ok := m[key]
	if !ok {
		return empty, errCrawlNotFound
	}

	tVal, ok := val.(T)
	if !ok {
		return empty, fmt.Errorf("could not convert value to provided type %T", empty)
	}

	return tVal, nil
}

// mapTypeToMapInterface will convert any provided map with a string as key
// into a map of type map[string]interface{}.
func mapTypeToMapInterface[T any](m map[string]T) map[string]interface{} {

	response := make(map[string]interface{})

	for key, value := range m {
		response[key] = value
	}
	return response
}
