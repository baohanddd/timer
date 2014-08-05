package common

import "menteslibres.net/gosexy/redis"
import "log"

func RedisNew(host string, port uint) *redis.Client {
	var client *redis.Client

	client = redis.New()

	err := client.Connect(host, port)

	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
		return nil
	}

	log.Println("Connected to redis-server.")

	return client
}
