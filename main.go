package main

import (
	"encoding/json"
	"flag"
	"github.com/rustedturnip/go-csgo-item-parser/csgo"
	"os"

	"github.com/rustedturnip/go-csgo-item-parser/parser"
)

var (
	csgoItemsLocation   string
	csgoEnglishLocation string
	outputLocation      string
)

func init() {
	flag.StringVar(&csgoItemsLocation, "csgo-items", "/items_game.txt", "the path to the csgo_items.txt file")
	flag.StringVar(&csgoEnglishLocation, "csgo-english", "/csgo_english.txt", "the path to the csgo_english.txt file")
	flag.StringVar(&outputLocation, "output", "/result.json", "the path to resulting json output file")
}

func main() {

	flag.Parse()

	// read data
	languageData, err := parser.Parse(csgoEnglishLocation)
	if err != nil {
		panic(err)
	}

	itemData, err := parser.Parse(csgoItemsLocation)
	if err != nil {
		panic(err)
	}

	// parse data
	allItems, err := csgo.New(languageData, itemData)
	if err != nil {
		panic(err)
	}

	// output data
	fo, err := os.Create(outputLocation)
	if err != nil {
		panic(err)
	}

	defer fo.Close()

	encoder := json.NewEncoder(fo)

	// set pretty-printing
	encoder.SetIndent("", "    ")

	err = encoder.Encode(allItems)
	if err != nil {
		panic(err)
	}
}
