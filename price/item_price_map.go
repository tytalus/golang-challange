package price

import "sync"

type ItemPriceMap struct {
	sync.RWMutex
	internal map[string]cachedPrice
}

func NewItemPricesMap() *ItemPriceMap {
	return &ItemPriceMap{
		internal: make(map[string]cachedPrice),
	}
}

func (rm *ItemPriceMap) Load(key string) (value cachedPrice, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *ItemPriceMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *ItemPriceMap) Store(key string, value cachedPrice) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
