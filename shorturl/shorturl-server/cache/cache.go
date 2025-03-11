package cache

const DefaultTTL = 30 * 86400

type KVCache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl int) error
	Destroy()
}
