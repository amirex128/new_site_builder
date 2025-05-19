package service

type IMemoryLoader interface {
	GetSampleFromCache(id string) (interface{}, bool)
}
