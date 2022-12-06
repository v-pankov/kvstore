package memcached

import "github.com/vdrpkv/kvstore/internal/pkg/memcached/core"

func mapCoreItems(coreItems []core.Item) []Item {
	items := make([]Item, len(coreItems))
	for i := 0; i < len(items); i++ {
		items[0] = Item{
			Key:   coreItems[i].Key,
			Flags: coreItems[i].Flags,
			Value: string(coreItems[i].Value),
		}
	}
	return items
}
