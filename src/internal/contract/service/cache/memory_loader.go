package cache

type IMemoryLoader interface {
	GetSampleFromCache(id string) (interface{}, bool)
}
