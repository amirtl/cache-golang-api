package database

import (
	"github.com/go-redis/redis"
	"log"
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
	err := db.Client.SAdd(key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Get(key string) []string {
	val, err := db.Client.SMembers(key).Result()
	if err != nil {
		log.Print("SMembers FAILED.", err)
	}
	return val
}
