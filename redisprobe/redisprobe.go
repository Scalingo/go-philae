package redisprobe

import (
	"context"
	"net/url"

	redis "github.com/redis/go-redis/v9"
	errgo "gopkg.in/errgo.v1"
)

type RedisProbe struct {
	name     string
	password string
	host     string
}

func NewRedisProbeFromURL(name, serviceUrl string) RedisProbe {
	url, _ := url.Parse(serviceUrl)
	password := ""

	if url.User != nil {
		password, _ = url.User.Password()
	}

	return NewRedisProbe(name, url.Host, password)
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

func (p RedisProbe) Check(ctx context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     p.host,
		Password: p.password,
		DB:       0,
	})
	defer client.Close()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return errgo.Notef(err, "unable to contact Redis host")
	}

	return nil
}
