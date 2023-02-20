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

// paintkit represents the image details of a skin, i.e. the available float
// range the skin can be in. Every entities.Skin has an associated paintkit.
type paintkit struct {
	Id          string
	Name        string
	Description string
	Rarity      string
	MinFloat    float64
	MaxFloat    float64
}

// mapToPaintkit converts the provided map into a paintkit providing
// all required parameters are present and of the correct type.
func mapToPaintkit(data map[string]interface{}, language *language) (*paintkit, error) {

	response := &paintkit{
		MinFloat: defaultMinFloat,
		MaxFloat: defaultMaxFloat,
	}

	// get Name
	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.New("Id (name) missing from paintkit")
	} else {
		response.Id = val
	}

	// get language Name Id
	if val, ok := data["description_tag"].(string); ok {
		name, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.Name = name
	}

	// get language Description Id
	if val, ok := data["description_string"].(string); ok {
		description, err := language.lookup(val)
		if err != nil {
			return nil, err
		}

		response.Description = description
	}

	// get min float
	if val, ok := data["wear_remap_min"].(string); ok {
		if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
			response.MinFloat = valFloat
		} else {
			return nil, errors.New("paintkit has non-float min float value (wear_remap_min)")
		}
	}

	// get max float
	if val, ok := data["wear_remap_max"].(string); ok {
		if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
			response.MaxFloat = valFloat
		} else {
			return nil, errors.New("paintkit has non-float max float value (wear_remap_max)")
		}
	}

	return response, nil
}

// getPaintkits gathers all Paintkits in the provided items data and returns them
// as map[paintkitId]paintkit.
func (c *csgoItems) getPaintkits() (map[string]*paintkit, error) {

	response := map[string]*paintkit{
		"vanilla": {
			Id:       "vanilla",
			MinFloat: defaultMinFloat,
			MaxFloat: defaultMaxFloat,
		},
	}

	rarities, err := crawlToType[map[string]interface{}](c.items, "paint_kits_rarity")
	if err != nil {
		return nil, fmt.Errorf("unable to extract paint_kits_rarity: %s", err.Error())
	}

	kits, err := crawlToType[map[string]interface{}](c.items, "paint_kits")
	if err != nil {
		return nil, errors.New("unable to locate paint_kits in provided items") // TODO improve error
	}

	for _, kit := range kits {
		mKit, ok := kit.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected paintkit layout in paint_kits")
		}

		converted, err := mapToPaintkit(mKit, c.language)
		if err != nil {
			return nil, err
		}

		if rarity, ok := rarities[converted.Id].(string); ok {
			converted.Rarity = rarity
		}

		response[converted.Id] = converted
	}

	return response, nil
}
