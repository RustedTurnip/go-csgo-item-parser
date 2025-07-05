package csgo

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type Quality struct {
	Id       string `json:"id"`
	Index    int    `json:"value"`
	Weight   int    `json:"weight"`
	HexColor string `json:"hexColor"`
}

func mapToQuality(id string, data map[string]interface{}, language *language) (*Quality, error) {

	response := &Quality{
		Id: id,
	}

	if val, ok := data["value"].(string); ok {
		if valInt, err := strconv.Atoi(val); err == nil {
			response.Index = valInt
		} else {
			return nil, errors.Wrap(err, fmt.Sprintf("unexpected index (value) type: %s", val))
		}
	} else {
		return nil, fmt.Errorf("quality (%s) missing expected field \"value\"", response.Id)
	}

	if val, ok := data["weight"].(string); ok {
		if valInt, err := strconv.Atoi(val); err == nil {
			response.Weight = valInt
		} else {
			return nil, errors.Wrap(err, fmt.Sprintf("unexpected index (value) type: %s", val))
		}
	} else {
		return nil, fmt.Errorf("quality (%s) missing expected field \"weight\"", response.Id)
	}

	if val, ok := data["hexColor"].(string); ok {
		response.HexColor = val
	} else {
		return nil, fmt.Errorf("quality (%s) missing expected field \"hexColor\"", response.Id)
	}

	return response, nil
}

func (c *csgoItems) getQualities() (map[string]*Quality, error) {
	response := make(map[string]*Quality)

	qualities, err := crawlToType[map[string]interface{}](c.items, "qualities")
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate qualities amongst items")
	}

	for id, quality := range qualities {
		qualityData, ok := quality.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Quality data for %s is of unexpected type", id)
		}

		qualityMap, err := mapToQuality(id, qualityData, c.language)
		if err != nil {
			return nil, err
		}

		response[id] = qualityMap
	}

	return response, nil
}
