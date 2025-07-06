package csgo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type stickerVariant string

const (
	stickerVariantPaper      stickerVariant = "Paper"
	stickerVariantGlossy     stickerVariant = "Glossy"
	stickerVariantGlitter    stickerVariant = "Glitter"
	stickerVariantHolo       stickerVariant = "Holo"
	stickerVariantFoil       stickerVariant = "Foil"
	stickerVariantGold       stickerVariant = "Gold"
	stickerVariantLenticular stickerVariant = "Lenticular"
)

var (
	stickerVariantIdSuffixes = map[string]stickerVariant{
		"_paper":      stickerVariantPaper,
		"_glossy":     stickerVariantGlossy,
		"_glitter":    stickerVariantGlitter,
		"_holo":       stickerVariantHolo,
		"_foil":       stickerVariantFoil,
		"_gold":       stickerVariantGold,
		"_lenticular": stickerVariantLenticular,
	}

	stickerVariantNameSuffixes = map[string]stickerVariant{
		"(Glitter)":    stickerVariantGlitter,
		"(Holo)":       stickerVariantHolo,
		"(Foil)":       stickerVariantFoil,
		"(Gold)":       stickerVariantGold,
		"(Lenticular)": stickerVariantLenticular,
	}
)

type stickerSubtypeContainer struct {
	stickers map[string]*Stickerkit
	sprays   map[string]*Spraykit
	patches  map[string]*Patchkit
}

// Stickerkit represents a Stickerkit object from the items_game file.
type Stickerkit struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RarityId    string `json:"rarityId"`
	Variant     string `json:"variant"`
}

// mapToStickerkit converts the provided data map into a Stickerkit object.
func mapToStickerkit(index int, data map[string]interface{}, language *language) (*Stickerkit, error) {

	response := &Stickerkit{
		Index:   index,
		Variant: "Paper",
	}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "Id (name) missing from Stickerkit")
	} else {
		response.Id = val
	}

	// identify sticker variant from id
	for suffix, variant := range stickerVariantIdSuffixes {
		if !strings.HasSuffix(response.Id, suffix) {
			continue
		}

		response.Variant = string(variant)
	}

	// get language Name
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("item_name missing from Stickerkit (%s)", response.Id))
	} else {

		lang, _ := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("language lookup of item_name for Stickerkit failed for key %s", val))
		}

		response.Name = lang
	}

	// identify sticker variant from descriptive name (which takes precedence)
	for suffix, variant := range stickerVariantNameSuffixes {
		if !strings.HasSuffix(response.Name, suffix) {
			continue
		}

		response.Variant = string(variant)
	}

	// get language Description
	if val, err := crawlToType[string](data, "description_string"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("description_string missing from Stickerkit (%s)", response.Id))
	} else {
		lang, _ := language.lookup(val)
		response.Description = lang
	}

	// get Rarity
	if val, err := crawlToType[string](data, "item_rarity"); err == nil {
		response.RarityId = val
	}

	return response, nil
}

// Spraykits are also stored as stickers in items_game, parsed into a separate container
type Spraykit struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RarityId    string `json:"rarityId"`
}

// mapToSpraykit converts the provided data map into a Spraykit object.
func mapToSpraykit(index int, data map[string]interface{}, language *language) (*Spraykit, error) {

	response := &Spraykit{
		Index: index,
	}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "Id (name) missing from SprayKit")
	} else {
		response.Id = val
	}

	// get language Name
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("item_name missing from SprayKit (%s)", response.Id))
	} else {

		lang, _ := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("language lookup of item_name for SprayKit failed for key %s", val))
		}

		response.Name = lang
	}

	// get language Description
	if val, err := crawlToType[string](data, "description_string"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("description_string missing from SprayKit (%s)", response.Id))
	} else {
		lang, _ := language.lookup(val[1:])
		response.Description = lang
	}

	// get Rarity
	if val, err := crawlToType[string](data, "item_rarity"); err == nil {
		response.RarityId = val
	}

	return response, nil
}

// Patchkits are also stored as stickers in items_game, parsed into a separate container
type Patchkit struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RarityId    string `json:"rarityId"`
}

// mapToPathkit converts the provided data map into a Patchkit object.
func mapToPathkit(index int, data map[string]interface{}, language *language) (*Patchkit, error) {

	response := &Patchkit{
		Index: index,
	}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "Id (name) missing from SprayKit")
	} else {
		response.Id = val
	}

	// get language Name
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("item_name missing from SprayKit (%s)", response.Id))
	} else {

		lang, _ := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("language lookup of item_name for SprayKit failed for key %s", val))
		}

		response.Name = lang
	}

	// get language Description
	if val, err := crawlToType[string](data, "description_string"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("description_string missing from SprayKit (%s)", response.Id))
	} else {
		lang, _ := language.lookup(val[1:])
		response.Description = lang
	}

	// get Rarity
	if val, err := crawlToType[string](data, "item_rarity"); err == nil {
		response.RarityId = val
	}

	return response, nil
}

// stickerSubtypeMapper chooses what function to use for sticker subtype classification
// return the converted sticker to caller
func stickerSubtypeMapper(index int, name string, data map[string]interface{}, language *language) (interface{}, error) {
	if strings.HasSuffix(name, "_graffiti") || strings.HasPrefix(name, "spray_") {
		return mapToSpraykit(index, data, language)
	} else if strings.Contains(name, "_teampatch_") || strings.HasPrefix(name, "patch_") {
		return mapToPathkit(index, data, language)
	} else {
		return mapToStickerkit(index, data, language)
	}
}

// getStickerkits retrieves all the stickerSubtypeContainer entries that are available
// in the provided items map and returns them.
func (c *csgoItems) getStickerkits() (*stickerSubtypeContainer, error) {

	response := &stickerSubtypeContainer{
		stickers: make(map[string]*Stickerkit),
		sprays:   make(map[string]*Spraykit),
		patches:  make(map[string]*Patchkit),
	}

	kits, err := crawlToType[map[string]interface{}](c.items, "sticker_kits")
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate sticker_kits in provided items")
	}

	for index, kit := range kits {

		iIndex, err := strconv.Atoi(index)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret Stickerkit index (%s) as int", iIndex))
		}

		mKit, ok := kit.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected Stickerkit layout in sticker_kits (at index %s)", index)
		}

		// get the name to see what subtype it is
		name, ok := mKit["name"].(string)
		if !ok {
			continue
		}

		// convert to a stickerSubtype
		converted, err := stickerSubtypeMapper(iIndex, name, mKit, c.language)
		if err != nil {
			return nil, err
		}

		switch t := converted.(type) {
		case *Stickerkit:
			response.stickers[t.Id] = t

		case *Spraykit:
			response.sprays[t.Id] = t

		case *Patchkit:
			response.patches[t.Id] = t
		}
	}

	return response, nil
}
