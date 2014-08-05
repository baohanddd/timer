package main

import (
	"log"
	"menteslibres.net/gosexy/redis"
)

func main() {
	var (
		client *redis.Client
		err    error
		s      string
	)

	host := "127.0.0.1"
	var port uint = 6379

	client = redis.New()

	err = client.Connect(host, port)

	if err != nil {
		log.Fatalf("Connect failed: %s\n", err.Error())
		return
	}

	log.Println("Connected to redis-server.")

	log.Printf("Sending PING...\n")
	s, err = client.Ping()

	if err != nil {
		log.Fatalf("Could not ping: %s\n", err.Error())
		return
	}

	log.Printf("Received %s!\n", s)

	client.Quit()
}
