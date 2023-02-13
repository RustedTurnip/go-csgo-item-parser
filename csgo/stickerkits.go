package csgo

import (
	"errors"
)

type stickerkit struct {
	id                    string
	languageNameId        string
	languageDescriptionId string
	rarity                string
}

func mapToStickerkit(data map[string]interface{}) (*stickerkit, error) {

	response := &stickerkit{}

	// get name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("id (name) missing from stickerkit") // TODO improve error
	} else {
		response.id = val
	}

	// get language name id
	if val, err := crawlToType[string](data, "item_name"); err != nil {
		return nil, errors.New("language name id (item_name) missing from stickerkit") // TODO improve error
	} else {
		response.languageNameId = val
	}

	// get language description id
	if val, err := crawlToType[string](data, "description_string"); err != nil {
		return nil, errors.New("language description id (description_string) missing from stickerkit") // TODO improve error
	} else {
		response.languageDescriptionId = val
	}

	// get rarity
	if val, err := crawlToType[string](data, "item_rarity"); err == nil {
		response.rarity = val
	}

	return response, nil
}

func getStickerkits(items map[string]interface{}) (map[string]*stickerkit, error) {

	response := make(map[string]*stickerkit)

	kits, err := crawlToType[map[string]interface{}](items, "sticker_kits")
	if err != nil {
		return nil, errors.New("unable to locate paint_kits in provided items") // TODO improve error
	}

	for _, kit := range kits {

		mKit, ok := kit.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected stickerkit layout in sticker_kits")
		}

		converted, err := mapToStickerkit(mKit)
		if err != nil {
			return nil, err
		}

		response[converted.id] = converted
	}

	return response, nil
}
