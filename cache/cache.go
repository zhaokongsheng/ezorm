package cache

type Cache interface {
	Get(key string, dest interface{}) error
}
