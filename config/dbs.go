package config

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func InitPostgres() *gorm.DB {
	var (
		dbUser, dbPass, dbName, dbHost, dbPort string
		err                                    error
	)

	if dbUser, err = GetEnv("DB_USER"); err != nil {
		log.Fatalf("DB_USER: %w", err)
	}
	if dbPass, err = GetEnv("DB_PASSWORD"); err != nil {
		log.Fatalf("DB_PASSWORD: %w", err)
	}
	if dbName, err = GetEnv("DB_NAME"); err != nil {
		log.Fatalf("DB_NAME: %w", err)
	}
	if dbHost, err = GetEnv("DB_HOST"); err != nil {
		log.Fatalf("DB_HOST: %w", err)
	}
	if dbPort, err = GetEnv("DB_PORT"); err != nil {
		log.Fatalf("DB_PORT: %w", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB_PORT: %w", err)
	}
	return db
}

func InitRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "redis:6379")
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Panicln("redis ping error:", err)
			}
			return err
		},
	}

	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("CONFIG", "SET", "maxmemory", "256mb")
	if err != nil {
		log.Panicln("redis config error:", err)
	}
	_, err = conn.Do("CONFIG", "SET", "maxmemory-policy", "allkeys-lru")
	if err != nil {
		log.Panicln("redis config error:", err)
	}

	return redisPool
}
