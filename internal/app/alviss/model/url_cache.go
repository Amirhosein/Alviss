package model

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

type URLRepo interface {
	Save(shortURL string, urlMapping URLMapping, expTime time.Duration) error
	Get(shortURL string) (URLMapping, error)
	Update(shortURL string, urlMapping URLMapping) error
}

type CacheURLRepo struct {
	URLDB
	RedisClient *redis.Client
}

func (s CacheURLRepo) Save(shortURL string, urlMapping URLMapping, expTime time.Duration) error {
	err := s.URLDB.Save(shortURL, urlMapping, expTime)
	if err != nil {
		return err
	}
	return nil
}

func (s CacheURLRepo) Get(shortURL string) (URLMapping, error) {
	var urlMapping URLMapping

	result, err := s.RedisClient.Get(shortURL).Result()
	urlMapping.UnmarshalBinary([]byte(result))
	if err == redis.Nil {
		log.Println("Redis cache miss")
		urlMapping, err = s.URLDB.Get(shortURL)
		if err != nil {
			return urlMapping, err
		}
		err = s.RedisClient.Set(shortURL, urlMapping, time.Hour*6).Err()
		return urlMapping, err
	} else if err != nil {
		return urlMapping, err
	} else {
		log.Println("Redis cache hit")
		return urlMapping, nil
	}
}

func (s CacheURLRepo) Update(shortURL string, urlMapping URLMapping) error {
	s.URLDB.Update(shortURL, urlMapping)
	err := s.RedisClient.Set(shortURL, urlMapping, 0).Err()
	return err
}
