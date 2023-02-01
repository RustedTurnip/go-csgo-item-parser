package main

import (
	"flag"
	"fmt"
)

var (
	csgoItemsLocation   string
	csgoEnglishLocation string
)

func init() {
	flag.StringVar(&csgoItemsLocation, "csgo-items", "", "the path to the csgo_items.txt file")
	flag.StringVar(&csgoEnglishLocation, "csgo-english", "/Users/samuel/Downloads/csgo_english.txt", "the path to the csgo_english.txt file")
}

func main() {
	flag.Parse()

	result, err := parse(csgoEnglishLocation)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(result["lang"].(map[string]interface{})["Tokens"].(map[string]interface{})))
}
