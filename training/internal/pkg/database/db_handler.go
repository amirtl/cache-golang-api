package database

import (
	"github.com/go-redis/redis"
	"log"
	"strings"
	"training/internal/pkg/config"
)

type Database struct {
	Client *redis.Client
}

func NewDatabase(conf *config.Base) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.DataBaseUrl,
		Password: conf.DataBasePassword,
		DB:       conf.DB,
	})
	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}
	return &Database{
		Client: client,
	}, nil
}

func (db *Database) Add(key string, value []string) error {
	for _, word := range value {
		val, err := db.Client.GetRange(key, 0, -1).Result()
		if val == "" {
			err = db.Client.Append(key, word).Err()
		} else {
			err = db.Client.Append(key, ","+word).Err()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Get(key string) []string {
	val, err := db.Client.GetRange(key, 0, -1).Result()
	if err != nil {
		log.Print("LRANGE FAILED.", err)
	}
	ValArray := strings.Split(val, ",")
	return ValArray
}
