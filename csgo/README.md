# csgo

This file represents notes for interpreting the `items_game.txt` file.

## items_game file structure

```
items_game
|
|_ rarities
   |
   |_ rarity_key [variable]
      |
      |_ value
      |_ loc_key
      |_ loc_weapon_key
      |_ loc_weapon_character
|
|_ qualities
   |
   |_ quality_key [variable]
      |
      |_ value
      |_ hexColor
|
|_ colors
   |
   |_ color_key [variable]
      |_ color_name
      |_ hex_color
|
|_ graffiti_tints
   |
   |_ graffiti_tint_key [variable]
      |
      |_ id
      |_ hex_color
|
|_ alternate_icons2
   |
   |_ weapon_icons
      |
      |_ weapon_icon_key [variable]
         |
         |_ icon_path
   |
   |_ casket_icons
      |
      |_ casket_icon_key [variable]
         |
         |_ icon_path
|
|_ prefabs
   |
   |_ prefab_key [variable]
|
|_ items
   |
   |_ item_key [variable]
      |
      |_ name
|
|_ music_definitions
   |
   |_ music_definition_key [variable]
      |
      |_ name
      |_ loc_name
      |_ loc_description
...
```


### items breakdown


**weapons**

- all weapons have a prefab, so to detect, simply lookup prefab map of
  prefabs we are interested in.


**collectible coins**

- `item_name` begins with `#CSGO_CollectibleCoin_`
- also have prefabs of `pickem_trophy`, `majors_trophy`, and
  `collectible_untradable_coin`


**season pass**

- `item_name` begins with `#CSGO_Ticket_`
- has prefab of `season_pass`

**collectible**

- `item_name` begins with `#CSGO_Collectible_`


**csgo tool**

- `item_name` begines with `#CSGO_Tool_` or `#CSGO_tool_`


**case/capsule key**

- has `prefab` of `valve weapon_case_key` - the prefab is the second part, or
  a prefab of just `weapon_case_key`


**campaign**

- `item_name` starts with #csgo_campaign
- prefab of `valve campaign_prefab` or `campaign_prefab`


**game license**

- is a coin that so far cannot be categorised


**weapon crate**

- has prefab of `weapon_case` or `weapon_case_souvenirpkg`
- `item_name` begins with `CSGO_crate`
- has `tags > ItemSet > tag_value` of a weapon set


**sticker capsule**

- has prefab of `sticker_capsule`
- examples without prefab:
    - `crate_signature_pack_stockh2021_group_finalists`
    - `CSGO_crate_signature_pack_atlanta2017_astr`
- stickers located through `attributes > set supply crate series > value`
  which points to a revolving_loot_list
- NOTE: their `item_name`s begin with `#CSGO_crate` and they link out to a
  `revolving_loot_list`. Maybve this can be used to identify `#CSGO_crate`
  types. IF THEY DON'T LINK OUT TO A REVOLVING_LOOT_LIST, THEN THEY HAVE A
  LOOT_LIST_NAME that can be used.


**operator_dossier**

- Example: `character_operator_dossier_op09_ancient`
- `loot_list_name` filed points directly to a `client_loot_list`


## prefabs breakdown

All items (except for `default`) have a prefab.

**weapon_case_base**

Can be:
- Sticker Pack (Capsule)
- Special Weapon Case (see `crate_xray_p250`)
- Pins capsule
- Music Kit capsule
- Operator Dossier
- Weapon Set (as item, beginning with `#CSGO_set_`), these can be ignored



