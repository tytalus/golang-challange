package price

import "sync"

type ItemPriceMap struct {
	mutex    sync.RWMutex
	internal map[string]cachedPrice
}

func NewItemPricesMap() *ItemPriceMap {
	return &ItemPriceMap{
		internal: make(map[string]cachedPrice),
	}
}

func (rm *ItemPriceMap) Load(key string) (value cachedPrice, ok bool) {
	rm.mutex.RLock()
	result, ok := rm.internal[key]
	rm.mutex.RUnlock()
	return result, ok
}

func (rm *ItemPriceMap) Delete(key string) {
	rm.mutex.Lock()
	delete(rm.internal, key)
	rm.mutex.Unlock()
}

func (rm *ItemPriceMap) Store(key string, value cachedPrice) {
	rm.mutex.Lock()
	rm.internal[key] = value
	rm.mutex.Unlock()
}
