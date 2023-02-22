package csgo

import (
	"fmt"
	"github.com/pkg/errors"
)

// stickerkit represents a stickerkit object from the items_game file.
type stickerkit struct {
	Id          string
	Name        string
	Description string
	RarityId    string
}

// mapToStickerkit converts the provided data map into a stickerkit object.
func mapToStickerkit(data map[string]interface{}, language *language) (*stickerkit, error) {

	response := &stickerkit{}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "Id (name) missing from stickerkit")
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("item_name missing from stickerkit (%s)", response.Id))
	} else {

		lang, _ := language.lookup(val)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("language lookup of item_name for stickerkit failed for key %s", val))
		}

		response.Name = lang
	}

	// get language Description Id
	if val, err := crawlToType[string](data, "description_string"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("description_string missing from stickerkit (%s)", response.Id))
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

// getStickerkits retrieves all the Stickerkits available in the provided items map
// and returns them in the format map[stickerkitId]stickerkit.
func (c *csgoItems) getStickerkits() (map[string]*stickerkit, error) {

	response := make(map[string]*stickerkit)

	kits, err := crawlToType[map[string]interface{}](c.items, "sticker_kits")
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate sticker_kits in provided items")
	}

	for index, kit := range kits {

		mKit, ok := kit.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected stickerkit layout in sticker_kits (at index %s)", index)
		}

		converted, err := mapToStickerkit(mKit, c.language)
		if err != nil {
			return nil, err
		}

		response[converted.Id] = converted
	}

	return response, nil
}
