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
	prices             *ItemPriceMap
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             NewItemPricesMap(),
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	priceFromCache, ok := c.prices.Load(itemCode)
	if ok {
		if time.Since(priceFromCache.RetrievedTime) < c.maxAge {
			return priceFromCache.Price, nil
		}
	}
	price, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v", err.Error())
	}
	c.prices.Store(itemCode, cachedPrice{Price: price, RetrievedTime: time.Now()})
	return price, nil
}

type channelResult struct {
	Price    float64
	Position int
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := make([]float64, len(itemCodes))
	resultChannel := make(chan channelResult)
	errorChannel := make(chan error)

	for i, itemCode := range itemCodes {
		go func(itemCode string, position int) {
			price, err := c.GetPriceFor(itemCode)
			if err != nil {
				errorChannel <- err
			}
			resultChannel <- channelResult{
				Price:    price,
				Position: position,
			}
		}(itemCode, i)
	}
	for range itemCodes {
		select {
		case result := <-resultChannel:
			results[result.Position] = result.Price
		case err := <-errorChannel:
			return []float64{}, err
		}
	}
	return results, nil
}
