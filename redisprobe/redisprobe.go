package redisprobe

import (
	errgo "gopkg.in/errgo.v1"
	redis "gopkg.in/redis.v3"
)

type RedisProbe struct {
	name     string
	password string
	host     string
}

func NewRedisProbe(name, host, password string) RedisProbe {
	return RedisProbe{
		name:     name,
		password: password,
		host:     host,
	}
}

func (p RedisProbe) Name() string {
	return p.name
}

func (p RedisProbe) Check() error {
	client := redis.NewClient(&redis.Options{
		Addr:     p.host,
		Password: "",
		DB:       0,
	})
	defer client.Close()

	_, err := client.Ping().Result()

	if err != nil {
		return errgo.Notef(err, "Unable to contact host")
	}

	return nil
}
