package main

import (
	"fmt"
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func main() {
	cfg, err := NewConfig()
	checkError(err)
	fmt.Println(cfg.Namespace)

	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port))
		},
	}

	var enqueuer = work.NewEnqueuer(cfg.Namespace, redisPool)

	emailAddress := []string{"mashelmi@yopmail.com"}

	for _, v := range emailAddress {
		_, err = enqueuer.Enqueue("send_welcome_email", work.Q{"email_address": v, "user_id": 4})
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Config struct {
	Namespace string
	Redis     RedisConfig
}

type RedisConfig struct {
	Host string
	Port string
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	checkError(err)
	err = viper.Unmarshal(&cfg)
	checkError(err)

	return cfg, nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
