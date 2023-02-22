# go-csgo-item-parser

`go-csgo-item-parser` can be used to extract CSGO entities from the items_game.txt and `csgo_<language>.txt`
files of the game.


## Support

Currently, the parser supports extraction of the following item types:

- paint kits
- sticker kits
- sticker capsules
- skinnable weapons (guns)
- skinnable weapons (knives)
- skinnable gloves
- weapon sets
- weapon crates


## Usage

The program currently works by accepting the file locations of the aforementioned `itmes_game.txt` and
`csgo_<language>.txt` files and outputs the transformed data into the provided output file location:

- `--csgo-items`: `items_game.txt` file location
- `--csgo-language`: `csgo_<language>.txt` file location
- `--output`: output file location

**Example**

```bash
go-csgo-item-parser --csgo-items=/path/to/items_game.txt --csgo-language=/path/to/csgo_english.txt --output=/path/to/result.json
```

The output file will contain the currently supported entities in json format.

