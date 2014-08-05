package common

import "github.com/fzzy/radix/redis"
import "time"
import "log"

func RedisClient(host string, port string) *redis.Client {
	c, err := redis.DialTimeout("tcp", host+":"+port, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	r := c.Cmd("select", 2)
	if r.Err != nil {
		log.Fatal(err)
	}

	return c
}

func RedisNew(host string, port string) *redis.Client {
	var client *redis.Client

	client = redis.New()

	err = client.Connect(host, port)

	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
		return
	}

	log.Println("Connected to redis-server.")

	return client

	// client.Quit()
}
