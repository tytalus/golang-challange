package price

import (
	"fmt"
	"time"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}

type cachedPrice struct {
	Price         float64
	RetrievedTime time.Time
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]cachedPrice
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]cachedPrice{},
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	priceFromCache, ok := c.prices[itemCode]
	if ok {
		if time.Since(priceFromCache.RetrievedTime) < c.maxAge {
			return priceFromCache.Price, nil
		}
	}
	price, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v", err.Error())
	}
	c.prices[itemCode] = cachedPrice{Price: price, RetrievedTime: time.Now()}
	return price, nil
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := []float64{}
	for _, itemCode := range itemCodes {
		// TODO: parallelize this, it can be optimized to not make the calls to the external service sequentially
		price, err := c.GetPriceFor(itemCode)
		if err != nil {
			return []float64{}, err
		}
		results = append(results, price)
	}
	return results, nil
}
