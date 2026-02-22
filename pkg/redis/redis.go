package redis

import (
	"log"

	"github.com/go-redis/redis"
)

type Client struct {
	Conn *redis.Client
}

func New(address, password string, db int) *Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	if _, err := conn.Ping().Result(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	}

	return &Client{Conn: conn}
}
