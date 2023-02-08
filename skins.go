package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type wear struct {
	name     string
	minValue float64
	maxValue float64
}

type skinQuality string

var (
	skinQualityNormal   skinQuality = ""
	skinQualityStattrak skinQuality = "StatTrak™"
	skinQualitySouvenir skinQuality = "Souvenir"

	prefabToQuality = map[string]skinQuality{
		"weapon_case_souvenirpkg": skinQualitySouvenir,
		"weapon_case":             skinQualityStattrak,
	}

	wears = []*wear{
		{
			name:     "Factory New",
			minValue: 0,
			maxValue: 0.07,
		},
		{
			name:     "Minimal Wear",
			minValue: 0.07,
			maxValue: 0.15,
		},
		{
			name:     "Field-Tested",
			minValue: 0.15,
			maxValue: 0.37,
		},
		{
			name:     "Well-Worn",
			minValue: 0.37,
			maxValue: 0.44,
		},
		{
			name:     "Battle-Scarred",
			minValue: 0.44,
			maxValue: 1,
		},
	}

	weaponNames = map[string]string{
		// Pistols
		"weapon_hkp2000":      "P2000",
		"weapon_usp_silencer": "USP-S",
		"weapon_glock":        "Glock-18",
		"weapon_elite":        "Dual Berettas",
		"weapon_p250":         "P250",
		"weapon_cz75a":        "CZ75-Auto",
		"weapon_fiveseven":    "Five-SeveN",
		"weapon_tec9":         "Tec-9",
		"weapon_deagle":       "Desert Eagle",
		"weapon_revolver":     "R8 Revolver",

		// Heavy
		"weapon_nova":     "Nova",
		"weapon_xm1014":   "XM1014",
		"weapon_mag7":     "MAG-7",
		"weapon_sawedoff": "Sawed-Off",
		"weapon_m249":     "M249",
		"weapon_negev":    "Negev",

		// Subs
		"weapon_mp9":   "MP9",
		"weapon_mac10": "MAC-10",
		"weapon_mp5sd": "MP5-SD",
		"weapon_mp7":   "MP7",
		"weapon_ump45": "UMP-45",
		"weapon_p90":   "P90",
		"weapon_bizon": "PP-19 Bizon",

		// Rifles
		"weapon_famas":         "Famas",
		"weapon_m4a1":          "M4A4",
		"weapon_m4a1_silencer": "M4A1-S",
		"weapon_galilar":       "Galil AR",
		"weapon_ak47":          "AK-47",
		"weapon_ssg08":         "SSG 08",
		"weapon_aug":           "AUG",
		"weapon_sg556":         "SG 553",
		"weapon_awp":           "AWP",
		"weapon_scar20":        "SCAR-20",
		"weapon_g3sg1":         "G3SG1",

		// Knives
		"weapon_bayonet":               "★ Bayonet",
		"weapon_knife_survival_bowie":  "★ Bowie Knife",
		"weapon_knife_butterfly":       "★ Butterfly Knife",
		"weapon_knife_css":             "★ Classic Knife",
		"weapon_knife_falchion":        "★ Falchion Knife",
		"weapon_knife_flip":            "★ Flip Knife",
		"weapon_knife_gut":             "★ Gut Knife",
		"weapon_knife_tactical":        "★ Huntsman Knife",
		"weapon_knife_karambit":        "★ Karambit",
		"weapon_knife_m9_bayonet":      "★ M9 Bayonet",
		"weapon_knife_gypsy_jackknife": "★ Navaja Knife",
		"weapon_knife_outdoor":         "★ Nomad Knife",
		"weapon_knife_cord":            "★ Paracord Knife",
		"weapon_knife_push":            "★ Shadow Daggers",
		"weapon_knife_skeleton":        "★ Skeleton Knife",
		"weapon_knife_stiletto":        "★ Stiletto Knife",
		"weapon_knife_canis":           "★ Survival Knife",
		"weapon_knife_widowmaker":      "★ Talon Knife",
		"weapon_knife_ursus":           "★ Ursus Knife",

		// Gloves
		"studded_brokenfang_gloves_operation10": "★ Broken Fang Gloves",
		"studded_bloodhound_gloves_bloodhound":  "★ Bloodhound Gloves",
		"sporty_gloves_sporty":                  "★ Sport Gloves",
		"slick_gloves_slick":                    "★ Driver Gloves",
		"leather_handwraps_handwrap":            "★ Hand Wraps",
		"motorcycle_gloves_motorcycle":          "★ Moto Gloves",
		"specialist_gloves_specialist":          "★ Specialist Gloves",
		"studded_hydra_gloves_bloodhound_hydra": "★ Hydra Gloves",
	}

	reItemPath = regexp.MustCompile("^econ/default_generated/([a-zA-Z0-9_-]+)(_light|_medium|_heavy)$")
)

// getSkins provides a list of all skins from the provided language and items
// maps.
func getSkins(language, items map[string]interface{}) ([]*skin, error) {

	// extract skin/weapon mapping
	skinsToWeapons, err := extractWeaponSkins(items)
	if err != nil {
		return nil, err
	}

	// get available quality for skins
	qualitites := getSkinQualities(items)

	// make tokens lowercase (as the case doesn't always match between language and items)
	tokens := mapToLower(language)

	kits, ok := items["paint_kits"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to locate paint_kits in items data")
	}

	skins := make([]*skin, 0)

	for _, v := range kits {

		skinKey := v.(map[string]interface{})["name"].(string)

		langTagRef, ok := v.(map[string]interface{})["description_tag"].(string)
		if !ok {
			fmt.Printf("gotcha language ref: %s\n", skinKey)
			continue
		}

		langTagRef = strings.TrimPrefix(langTagRef, "#")

		langValue, ok := tokens[strings.ToLower(langTagRef)]
		if !ok {
			fmt.Printf("gotcha language value: %s\n", langTagRef)
		}

		for weapon, _ := range skinsToWeapons[skinKey] {

			for _, wear := range wears {

				// default float values (see: https://www.reddit.com/r/GlobalOffensiveTrade/comments/5fj8vm/comment/dakumlf/?utm_source=share&utm_medium=web2x&context=3)
				min := float64(0.06)
				max := float64(0.8)

				if strF, ok := v.(map[string]interface{})["wear_remap_min"].(string); ok {
					min, _ = strconv.ParseFloat(strF, 64)
				}

				if strF, ok := v.(map[string]interface{})["wear_remap_max"].(string); ok {
					max, _ = strconv.ParseFloat(strF, 64)
				}

				if max > wear.minValue && min <= wear.maxValue {

					// build normal version of skin
					skins = append(skins, &skin{
						MarketHashName: buildSkinMarketHashName(skinQualityNormal, weaponNames[weapon], langValue.(string), wear.name),
						Wear:           wear.name,
						MaxFloat:       max,
						MinFloat:       min,
					})

					// build special quality of skin if exists
					if quality, ok := qualitites[fmt.Sprintf("[%s]%s", skinKey, weapon)]; ok {
						skins = append(skins, &skin{
							MarketHashName: buildSkinMarketHashName(quality, weaponNames[weapon], langValue.(string), wear.name),
							Wear:           wear.name,
							MaxFloat:       max,
							MinFloat:       min,
						})
					}
				}
			}
		}
	}

	return skins, nil
}

// extractWeaponSkins builds a map of all skins and the weapons that they are
// available on.
func extractWeaponSkins(items map[string]interface{}) (map[string]map[string]interface{}, error) {

	skinsToWeapons := make(map[string]map[string]interface{})

	for _, weaponSkin := range items["alternate_icons2"].(map[string]interface{})["weapon_icons"].(map[string]interface{}) {
		path := weaponSkin.(map[string]interface{})["icon_path"].(string)

		weapon, skin, err := splitIconPath(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, path)
		}

		if _, ok := skinsToWeapons[skin]; !ok {
			skinsToWeapons[skin] = make(map[string]interface{})
		}

		skinsToWeapons[skin][weapon] = struct{}{}
	}

	return skinsToWeapons, nil
}

// splitIconPath will attempt to divide the provided "icon path" string into
// a weapon id and a paintkit id.
//
// An error is returned when the string cannot be matched to the expected path
// format.
func splitIconPath(path string) (string, string, error) {

	result := reItemPath.FindStringSubmatch(path)
	if len(result) != 3 {
		return "", "", errors.New("invalid path found")
	}

	fileParts := strings.Split(result[1], "_")
	for i := 0; i < len(fileParts); i++ {

		weapon := strings.Join(fileParts[:i], "_")

		if _, ok := weaponNames[weapon]; !ok {
			continue
		}

		return weapon, strings.Join(fileParts[i:], "_"), nil
	}

	return "", "", errors.New("unable to distinguish weapon and skin")
}

// getSkinQualities produces a mapping of "[skin]weapon" to available "special"
// quality if one is possible, e.g. Souvenir or StatTrak.
func getSkinQualities(items map[string]interface{}) map[string]skinQuality {

	setQualities := make(map[string]skinQuality)

	// loop through all items
	for _, itemData := range items["items"].(map[string]interface{}) {

		dataMap := itemData.(map[string]interface{})

		// if item is not a case prefab
		prefab, ok := dataMap["prefab"]
		if !ok {
			continue
		}

		// see if crate/set supports special skin qualities
		quality, ok := prefabToQuality[prefab.(string)]
		if !ok {
			continue
		}

		// grab item set name for souvenir
		if tags, ok := dataMap["tags"].(map[string]interface{}); ok {
			if itemSet, ok := tags["ItemSet"].(map[string]interface{}); ok {
				if setName, ok := itemSet["tag_value"].(string); ok {
					setQualities[setName] = quality
				}
			}
		}
	}

	// loop through all sets found and retrieve items
	skinQualities := make(map[string]skinQuality)
	for setName, quality := range setQualities {
		for item, _ := range items["item_sets"].(map[string]interface{})[setName].(map[string]interface{})["items"].(map[string]interface{}) {
			// TODO item is in format [gs_awp_hydra]weapon_awp - this needs to be considered
			skinQualities[item] = quality
		}
	}

	return skinQualities
}

// buildSkinMarketHashName takes the required attributes for a skin's market hash
// name and formats it into the uniquely identifiable market hash name.
func buildSkinMarketHashName(quality skinQuality, gun, skinName, wear string) string {

	if quality != skinQualityNormal {
		if strings.HasPrefix(gun, "★") {
			gun = fmt.Sprintf("★ %s %s", quality, strings.TrimPrefix(gun, "★ "))
		} else {
			gun = fmt.Sprintf("%s %s", quality, gun)
		}
	}

	if skinName == "" {
		return gun
	}

	return fmt.Sprintf("%s | %s (%s)", gun, skinName, wear)
}

// mapToLower takes a map with a string key type, and returns a copy with
// all keys converted to lowercase.
//
// Note: if there is more than one key with the same string (but with different
// cases) then only one of the values will be kept.
func mapToLower[T any](m map[string]T) map[string]T {

	nm := make(map[string]T)

	for k, v := range m {
		nm[strings.ToLower(k)] = v
	}

	return nm
}
