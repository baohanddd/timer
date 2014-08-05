package common

import "menteslibres.net/gosexy/redis"
import "log"

func RedisNew(host string, port uint) *redis.Client {
	var client *redis.Client

	client = redis.New()

	err := client.Connect(host, port)

	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
	}

	_, err = client.Select(2)

	if err != nil {
		log.Fatal("Can not select db: 2, due to ", err)
	}

	log.Println("Connected to redis-server.")

	return client
}
