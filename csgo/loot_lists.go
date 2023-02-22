package csgo

import (
	"fmt"
	"strings"
)

// revolvingLootLists represents the key value store of revolvingLootList ids (indexes)
// to their corresponding client_loot_list ids.
type revolvingLootLists map[string]string

// getRevolvingLootLists pulls all revolving_llot_list index to client_loot_list id pairs
// out of the items_game data and returns them as a revolvingLootLists.
func (c *csgoItems) getRevolvingLootLists() (revolvingLootLists, error) {

	response := make(revolvingLootLists)

	lootLists, err := crawlToType[map[string]interface{}](c.items, "revolving_loot_lists")
	if err != nil {
		return nil, err
	}

	for key, val := range lootLists {
		response[key] = val.(string)
	}

	return response, nil
}

// clientLootListItemType represents the underlying type of object stored within a
// client_loot_list.
type clientLootListItemType int

const (
	clientLootListItemTypeUnknown clientLootListItemType = iota
	clientLootListItemTypeSubList
	clientLootListItemTypeSticker
)

// clientLootListItems represents an item list for items within a client_loot_list,
// along with the items' type.
type clientLootListItems struct {
	listType clientLootListItemType
	items    []string
}

// clientLootList represents a flattened client_loot_list structure from the items_game
// file. The root of each client_loot_list object is available from the
// revolving_loot_list entities, and the items contained within clientLootList
// do not retain any subgroups of the client_loot_list.
type clientLootList struct {
	id        string
	listItems *clientLootListItems
}

// getClientLootLists retrieves all the client_loot_lists from the c.items
// map.
func (c *csgoItems) getClientLootLists() (map[string]*clientLootList, error) {

	response := make(map[string]*clientLootList)

	lootLists, err := crawlToType[map[string]interface{}](c.items, "client_loot_lists")
	if err != nil {
		return nil, err
	}

	// build map of clientLootList ids to clientLootList indexes
	for id, list := range lootLists {

		listMap, ok := list.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected client_loot_list Id format for revolving_loot_list index %s", id)
		}

		entry := &clientLootList{
			id:        id,
			listItems: &clientLootListItems{},
		}

		for itemName, _ := range listMap {

			// if list contains sublist
			if _, ok := lootLists[itemName]; ok {
				entry.listItems.listType = clientLootListItemTypeSubList
				entry.listItems.items = append(entry.listItems.items, itemName)
				continue
			}

			// if list contains stickers
			if strings.HasSuffix(itemName, "]sticker") {
				entry.listItems.listType = clientLootListItemTypeSticker
				itemName = strings.TrimPrefix(itemName, "[")
				itemName = strings.TrimSuffix(itemName, "]sticker")
				entry.listItems.items = append(entry.listItems.items, itemName)
				continue
			}

			continue
		}

		response[entry.id] = entry
	}

	return response, nil
}

// crawlClientLootLists will recursively traverse down through the lists that have
// clientLootListItemTypeSubList as the type to identify the root type of the list.
func crawlClientLootLists(listId string, clientLootLists map[string]*clientLootList) (clientLootListItemType, []string) {

	list, ok := clientLootLists[listId]
	if !ok {
		return clientLootListItemTypeUnknown, nil
	}

	switch list.listItems.listType {

	case clientLootListItemTypeSticker:
		return clientLootListItemTypeSticker, list.listItems.items

	case clientLootListItemTypeSubList:
		responseType := clientLootListItemTypeUnknown
		responseItems := make([]string, 0)

		for _, sublist := range list.listItems.items {
			subType, items := crawlClientLootLists(sublist, clientLootLists)
			responseType = subType
			responseItems = append(responseItems, items...)
		}

		return responseType, responseItems

	}

	return clientLootListItemTypeUnknown, nil
}
