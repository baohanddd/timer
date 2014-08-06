package common

import (
	"errors"
	"fmt"
	"log"
	"menteslibres.net/gosexy/redis"
	"strconv"
	"strings"
)

type RedisAddress struct {
	Host string
	Port uint
}

func (o *RedisAddress) SetPort(p string) {
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatalf("Can not convert port string to int: %s\n", p)
	}
	o.Port = uint(port)
}

func RedisNew(address string) *redis.Client {
	var (
		client *redis.Client
	)

	addr, err := parse(address)

	if err != nil {
		log.Fatal(err)
	}

	client = redis.New()

	err = client.Connect(addr.Host, addr.Port)

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

func parse(addr string) (*RedisAddress, error) {
	cfg := RedisAddress{}

	items := strings.Split(addr, ":")
	if len(items) != 2 {
		return nil, errors.New(fmt.Sprintf("Can not parse redis address: %s\n", addr))
	}

	cfg.Host = items[0]
	cfg.SetPort(items[1])

	return &cfg, nil
}
