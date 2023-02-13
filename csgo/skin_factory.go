package csgo

import (
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/rustedturnip/go-csgo-item-parser/entities"
)

type wear struct {
	name     string
	minValue float64
	maxValue float64
}

var (
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
)

func (c *csgo) getAllSkins() ([]*entities.Skin, error) {

	skins := make(map[string]*entities.Skin)

	// build crate skins
	for _, crate := range c.weaponCrates {

		set, ok := c.weaponSets[crate.weaponSetId]
		if !ok {
			continue
		}

		containerName, err := c.language.lookup(crate.languageNameId) // TODO give weaponCrate a getLanguageId function
		if err != nil {
			return nil, err
		}

		collectionName, err := c.language.lookup(set.languageNameId) // TODO give collection a getLanguageId function
		if err != nil {
			return nil, err
		}

		setSkins, err := c.weaponSetToSkins(crate.qualityCapability, set.items)
		if err != nil {
			return nil, err
		}

		for _, item := range setSkins {

			item.Collection = collectionName
			item.Containers = []string{containerName}

			if existing, ok := skins[item.MarketHashName]; ok {
				existing.Containers = append(existing.Containers, item.Containers...)
				continue
			}

			skins[item.MarketHashName] = item
		}
	}

	// build knife skins
	knifeSkins, err := c.weaponSetToSkins(qualityStatTrak, c.knifeSet)
	if err != nil {
		return nil, err
	}

	for _, knife := range knifeSkins {
		skins[knife.MarketHashName] = knife
	}

	// build glove skins
	gloveSkins, err := c.glovesSetToSkins(c.gloveSet)
	if err != nil {
		return nil, err
	}

	for _, gloves := range gloveSkins {
		skins[gloves.MarketHashName] = gloves
	}

	return maps.Values(skins), nil
}

// TODO comment
func (c *csgo) weaponSetToSkins(qualityCapability qualityCapability, set *itemPaintkitSet) ([]*entities.Skin, error) {

	skins := make([]*entities.Skin, 0)

	err := set.forEachItemPaintkit(func(itemId, paintkitId string) error {

		weapon, ok := c.weapons[itemId]
		if !ok {
			return fmt.Errorf("unexpected weapon id: %s", itemId)
		}

		paintkit, ok := c.paintkits[paintkitId]
		if !ok {
			return fmt.Errorf("unexpected paintkit id: %s", paintkitId)
		}

		newSkins, err := c.paintkitToSkins(qualityCapability, weapon, paintkit)
		if err != nil {
			return err
		}

		skins = append(skins, newSkins...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return skins, nil
}

// TODO comment
func (c *csgo) glovesSetToSkins(set *itemPaintkitSet) ([]*entities.Skin, error) {

	skins := make([]*entities.Skin, 0)

	err := set.forEachItemPaintkit(func(itemId, paintkitId string) error {

		gloves, ok := c.gloves[itemId]
		if !ok {
			return fmt.Errorf("unexpected weapon id: %s", itemId)
		}

		paintkit, ok := c.paintkits[paintkitId]
		if !ok {
			return fmt.Errorf("unexpected paintkit id: %s", paintkitId)
		}

		newSkins, err := c.paintkitToSkins(qualityNormal, gloves, paintkit)
		if err != nil {
			return err
		}

		skins = append(skins, newSkins...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return skins, nil
}

// TODO comment
func (c *csgo) paintkitToSkins(qualityCapability qualityCapability, item skinableItem, paintkit *paintkit) ([]*entities.Skin, error) {

	skins := make([]*entities.Skin, 0)

	weaponName, err := c.language.lookup(item.getLanguageNameId())
	if err != nil {
		return nil, err
	}

	paintkitName, err := c.language.lookup(paintkit.languageNameId)
	if err != nil {
		return nil, err
	}

	for _, wear := range getPaintkitAvailableWears(paintkit) {
		skins = append(skins, &entities.Skin{
			MarketHashName: buildSkinMarketHashName(item.getSpecial(), qualityCapability, paintkitName, weaponName, wear.name),
			// TODO do min/max floats
		})

		if qualityCapability != qualityNormal {
			skins = append(skins, &entities.Skin{
				MarketHashName: buildSkinMarketHashName(item.getSpecial(), qualityCapability, paintkitName, weaponName, wear.name),
				// TODO do min/max floats
			})
		}
	}

	return skins, nil
}

// TODO comment
func getPaintkitAvailableWears(paintkit *paintkit) []*wear {

	available := make([]*wear, 0)

	for _, wear := range wears {

		if paintkit.maxFloat > wear.minValue && paintkit.minFloat <= wear.maxValue {
			available = append(available, wear)
		}
	}

	return available
}

// buildSkinMarketHashName takes the required attributes for a skin's market hash
// name and formats it into the uniquely identifiable market hash name.
func buildSkinMarketHashName(isSpecial bool, quality qualityCapability, paintkitName, weaponName, wearName string) string {

	prefix := ""

	if isSpecial {
		prefix = "★ "
	}

	if quality != qualityNormal {
		prefix += fmt.Sprintf("%s ", string(quality))
	}

	return fmt.Sprintf("%s%s | %s (%s)", prefix, weaponName, paintkitName, wearName)
}
