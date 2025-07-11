package csgo

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type Musickit struct {
	Id          string `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func mapToMusickit(index int, data map[string]interface{}, language *language) (*Musickit, error) {
	response := &Musickit{
		Index: index,
	}

	if val, err := crawlToType[string](data, "name"); err != nil {
		return nil, errors.Wrap(err, "Id (name) missing from Musickit")
	} else {
		response.Id = val
	}

	if val, err := crawlToType[string](data, "loc_name"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("loc_name missing from Musickit (%s)", response.Id))
	} else {
		// english file omits the # at the beggining for chains
		lang, _ := language.lookup(val[1:])
		response.Name = lang
	}

	if val, err := crawlToType[string](data, "loc_description"); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("loc_description missing from Musickit (%s)", response.Id))
	} else {
		// english file omits the # at the beggining for chains
		lang, _ := language.lookup(val[1:])
		response.Description = lang
	}

	return response, nil
}

func (c *csgoItems) getMusickits() (map[string]*Musickit, error) {
	response := make(map[string]*Musickit)

	musickits, err := crawlToType[map[string]interface{}](c.items, "music_definitions")
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate music_definitions in provided items")
	}
	for index, music := range musickits {

		musicData, ok := music.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Musickit data for %s is of unexpected type", index)
		}

		iIndex, err := strconv.Atoi(index)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret Musickit index (%s) as int", index))
		}

		musicMap, err := mapToMusickit(iIndex, musicData, c.language)

		response[musicMap.Id] = musicMap

	}

	return response, nil
}
