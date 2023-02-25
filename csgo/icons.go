package csgo

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// getKnifeSet is used in the same way as getIconSet, but it supplements the set
// with additional vanilla items of the provided knifeIds.
func (c *csgoItems) getKnifeSet(knifeIds map[string]interface{}) (map[string][]string, error) {

	set, err := c.getIconSet(knifeIds)
	if err != nil {
		return nil, err
	}

	// add vanilla knives
	for knifeId := range knifeIds {
		set["vanilla"] = append(set["vanilla"], knifeId)
	}

	return set, nil
}

// getIconSet is used to extract the Weapon id-Paintkit id combinations from the
// alternate_icons2 list from items_game.txt. This can be used to extract items
// that do not appear in any sets defined elsewhere within the file.
func (c *csgoItems) getIconSet(itemIds map[string]interface{}) (map[string][]string, error) {

	response := make(map[string][]string)

	icons, err := crawlToType[map[string]interface{}](c.items, "alternate_icons2", "weapon_icons")
	if err != nil {
		return nil, fmt.Errorf("unable to locate weapon_icons: %s", err.Error())
	}

	for index, data := range icons {

		iconMap, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected alternate_icons2 format %s", index)
		}

		iconPath, err := crawlToType[string](iconMap, "icon_path")
		if err != nil {
			return nil, errors.Wrap(err, "couldn't crawl to path: icon_path")
		}

		// each Weapon skin appears in icons 3 times, by including only the "..._light"
		// ones we are eliminating duplicates.
		if !strings.HasSuffix(iconPath, "_light") {
			continue
		}

		targetId := findLongestIdMatch(itemIds, iconPath)
		if targetId == "" {
			continue
		}

		itemId, paintkitId, err := getItemPaintkitFromIconPath(targetId, iconPath)
		if err != nil {
			return nil, err
		}

		response[paintkitId] = append(response[paintkitId], itemId)
	}

	return response, nil
}

// findLongestIdMatch will take a provided icon path and locate the
// longest matching Id from ids within the path.
//
// e.g. with the ids: { test_weapon_knife, test_weapon_knife_karambit }
// and the path "icon/path/test_weapon_knife_karambit",
// test_weapon_knife_karambit will be returned.
func findLongestIdMatch[T any](ids map[string]T, path string) string {

	longest := ""

	for id, _ := range ids {

		if !strings.Contains(path, id) {
			continue
		}

		if len(id) > len(longest) {
			longest = id
		}
	}

	return longest
}

// getItemPaintkitFromIconPath will extract from the provided path, the
// Paintkit ID. itemID is required to distinguish the Paintkit from the
// Weapon.
func getItemPaintkitFromIconPath(itemId string, path string) (string, string, error) {

	pathTail := strings.TrimPrefix(path, "econ/default_generated/")
	components := strings.Split(pathTail, "_")

	for i := 0; i < len(components); i++ {

		iId := strings.Join(components[:i], "_")
		pkId := strings.Join(components[i:len(components)-1], "_") // drop last component (as light, medium or heavy)

		if iId != itemId {
			continue
		}

		return iId, pkId, nil
	}

	return "", "", fmt.Errorf("unable to derive Weapon and Paintkit from icon path: %s", path)
}
