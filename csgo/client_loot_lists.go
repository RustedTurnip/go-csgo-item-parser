package csgo

import (
	"fmt"
)

// clientLootList represents a flattened client_loot_list structure from the items_game
// file. The root of each client_loot_list object is available from the
// revolving_loot_list entities, and the items contained within clientLootList
// do not retain any subgroups of the client_loot_list.
type clientLootList struct {
	id    string
	index string
	items []string
}

// getClientLootLists retrieves all the client_loot_lists from the c.items
// map.
func (c *csgoItems) getClientLootLists() (map[string]*clientLootList, error) {

	response := make(map[string]*clientLootList)

	lootLists, err := crawlToType[map[string]interface{}](c.items, "client_loot_lists")
	if err != nil {
		return nil, err
	}

	revolvingLootLists, err := crawlToType[map[string]interface{}](c.items, "revolving_loot_lists")
	if err != nil {
		return nil, err
	}

	// build map of clientLootList ids to clientLootList indexes
	for index, clientLootListId := range revolvingLootLists {

		id, ok := clientLootListId.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected client_loot_list Id format for revolving_loot_list index %s", id)
		}

		response[index] = &clientLootList{
			id:    id,
			index: index,
			items: crawlClientLootLists(id, lootLists),
		}
	}

	return response, nil
}

// crawlClientLootLists will recursively traverse down through the groups of the provided
// list id and fetch all items contained within. Any subgroups will be flattened so the
// resulting list is flat.
func crawlClientLootLists(listId string, clientLootLists map[string]interface{}) []string {

	response := make([]string, 0)

	sublist, ok := clientLootLists[listId].(map[string]interface{})
	if !ok {
		return response
	}

	for sublistItem, _ := range sublist {

		if _, ok := clientLootLists[sublistItem]; ok {
			response = append(response, crawlClientLootLists(sublistItem, clientLootLists)...)
			continue
		}

		response = append(response, sublistItem)
	}

	return response
}
