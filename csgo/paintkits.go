package csgo

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	defaultMinFloat float64 = 0.06
	defaultMaxFloat float64 = 0.8
)

// TODO comment struct
type paintkit struct {
	id                    string
	languageNameId        string
	languageDescriptionId string
	rarityId              string
	minFloat              float64
	maxFloat              float64
}

// mapToPaintkit converts the provided map into a paintkit providing
// all required parameters are present and of the correct type.
func mapToPaintkit(data map[string]interface{}) (*paintkit, error) {

	response := &paintkit{
		minFloat: defaultMinFloat,
		maxFloat: defaultMaxFloat,
	}

	// get name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("id (name) missing from paintkit")
	} else {
		response.id = val
	}

	// get language name id
	if val, ok := data["description_tag"].(string); ok {
		response.languageNameId = val
	}

	// get language description id
	if val, ok := data["description_string"].(string); ok {
		response.languageDescriptionId = val
	}

	// get min float
	if val, ok := data["wear_remap_min"].(string); ok {
		if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
			response.minFloat = valFloat
		} else {
			return nil, errors.New("paintkit has non-float min float value (wear_remap_min)")
		}
	}

	// get max float
	if val, ok := data["wear_remap_max"].(string); ok {
		if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
			response.minFloat = valFloat
		} else {
			return nil, errors.New("paintkit has non-float max float value (wear_remap_max)")
		}
	}

	return response, nil
}

func getPaintkits(items map[string]interface{}) (map[string]*paintkit, error) {

	response := make(map[string]*paintkit)

	rarities, err := crawlToType[map[string]interface{}](items, "paint_kits_rarity")
	if err != nil {
		return nil, fmt.Errorf("unable to extract paint_kits_rarity: %s", err.Error())
	}

	kits, err := crawlToType[map[string]interface{}](items, "paint_kits")
	if err != nil {
		return nil, errors.New("unable to locate paint_kits in provided items") // TODO improve error
	}

	for _, kit := range kits {
		mKit, ok := kit.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected paintkit layout in paint_kits")
		}

		converted, err := mapToPaintkit(mKit)
		if err != nil {
			return nil, err
		}

		if rarity, ok := rarities[converted.id].(string); ok {
			converted.rarityId = rarity
		}

		response[converted.id] = converted
	}

	return response, nil
}
