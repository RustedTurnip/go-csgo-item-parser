package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/rustedturnip/go-csgo-item-parser/parser"
)

var (
	csgoItemsLocation   string
	csgoEnglishLocation string
)

func init() {
	flag.StringVar(&csgoItemsLocation, "csgo-items", "/Users/samuel/Downloads/items_game.txt", "the path to the csgo_items.txt file")
	flag.StringVar(&csgoEnglishLocation, "csgo-english", "/Users/samuel/Downloads/csgo_english.txt", "the path to the csgo_english.txt file")
}

func main() {
	flag.Parse()

	result, err := parser.Parse(csgoEnglishLocation)
	if err != nil {
		panic(err)
	}

	resultTwo, err := parser.Parse(csgoItemsLocation)
	if err != nil {
		panic(err)
	}

	allItems, err := getItems(result, resultTwo)
	if err != nil {
		panic(err)
	}

	for _, skin := range allItems.Skins {
		fmt.Println(skin.MarketHashName)
	}
}

// getItems retrieves all items from the provided items/language file and
// returns them as an items struct.
func getItems(languageData, itemData map[string]interface{}) (*items, error) {

	// check language base data exists
	lang, ok := languageData["lang"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to locate \"lang\" in provided languageData")
	}

	tokens, ok := lang["Tokens"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to locate \"lang/Tokens\" in provided languageData")
	}

	// check items base data exists
	fileItems, ok := itemData["items_game"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to locate \"items_game\" in provided itemData")
	}

	// extract skins
	skins, err := getSkins(tokens, fileItems)
	if err != nil {
		return nil, fmt.Errorf("unable to extract skins with error: %s", err.Error())
	}

	return &items{
		Skins: skins,
	}, nil
}
