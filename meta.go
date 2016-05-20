// Package imageEncrypt meta information
package imageEncrypt

import (
	"errors"
	"time"

	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

// Meta interface of meta information
type Meta interface {
	Save(metaImage MetaCuttedImage, condition ...interface{}) (interface{}, error)
	Get(condition ...interface{}) (MetaCuttedImage, error)
}

// MetaByRedis Use redis store the meta info
type MetaByRedis struct {
	pool *redis.Pool
}

// NewMetaByRedis constructor
func NewMetaByRedis(addr, pass string) *MetaByRedis {
	pool := newPool(addr, pass)
	return &MetaByRedis{pool}
}

func newPool(addr, pass string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     2,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				panic(err)
			}
			if pass != "" {
				if _, err = c.Do("AUTH", pass); err != nil {
					c.Close()
					panic(err)
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				panic(err)
			}
			return err
		},
	}
}

func (m *MetaByRedis) Save(metaImage MetaCuttedImage, condition ...interface{}) (interface{}, error) {
	data, err := json.Marshal(metaImage)
	if err != nil {
		return nil, err
	}
	return m.pool.Get().Do("SET", condition[0], data)
}

func (m *MetaByRedis) Get(condition ...interface{}) (MetaCuttedImage, error) {
	data, err := m.pool.Get().Do("GET", condition[0])
	metaImage := MetaCuttedImage{}
	if err != nil {
		return metaImage, err
	}
	if data == nil {
		return metaImage, errors.New("meta info is not exist!")
	}
	err = json.Unmarshal(data.([]byte), &metaImage)
	return metaImage, err
}
