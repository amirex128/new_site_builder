package service

type IMemoryLoaderService interface {
	GetSampleFromCache(id string) (interface{}, bool)
}
