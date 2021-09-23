package rcache

import (
	"encoding/json"
	"io"
	"log"

	cmap "github.com/MoSunDay/concurrent-map"
)

type cacheManager struct {
	data *cmap.ConcurrentMap
}

func NewCacheManager() *cacheManager {
	cm := &cacheManager{}
	cm.data = cmap.New(128)
	return cm
}

func (c *cacheManager) Get(key string) (ret string) {
	if value, ok := c.data.Get(key); ok {
		ret = value.(string)
	} else {
		ret = ""
		log.Printf("cacheManager get key: %s failed\n", key)
	}
	return ret
}

func (c *cacheManager) Set(key string, value string) error {
	c.data.Set(key, value)
	return nil
}

// Marshal serializes cache data
func (c *cacheManager) Marshal() ([]byte, error) {

	cacheKV := make(map[string]string, 20000)
	for _, key := range c.data.Keys() {
		value, result := c.data.Get(key)

		if result {
			cacheKV[key] = value.(string)
		}
	}
	dataBytes, err := json.Marshal(cacheKV)
	return dataBytes, err
}

// UnMarshal deserializes cache data
func (c *cacheManager) UnMarshal(serialized io.ReadCloser) error {
	var newData map[string]string
	if err := json.NewDecoder(serialized).Decode(&newData); err != nil {
		return err
	}

	for k, v := range newData {
		c.data.Set(k, v)
	}
	return nil
}
