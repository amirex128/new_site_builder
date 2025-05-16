package globalloader

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var dealProjectCache *cache.Cache

func DealProjectLoader() {
	dealProjectCache = cache.New(time.Minute, 2*time.Minute)

	dealProjectCache.Set("sample:1", "1", time.Minute)

}

func GetSampleFromCache(id string) (interface{}, bool) {
	if dealProjectCache == nil {
		return nil, false
	}
	return dealProjectCache.Get("sample:" + id)
}
