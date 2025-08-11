package config

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func createDatabaseIfNotExists(dbPass, dbHost, dbPort, dbName string) error {
	connStr := fmt.Sprintf("host=%s port=%s user=postgres password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbPass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(dbName)))
		if err != nil {
			return err
		}
		log.Println("Database created")
	} else {
		log.Println("Database already exists")
	}

	db.Close()

	connStr = fmt.Sprintf("host=%s port=%s user=postgres password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbPass, dbName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	migrator := MustGetNewMigrator(MigrationsFS, migrationsDir)
	err = migrator.ApplyMigrations(db)
	if err != nil {
		panic(err)
	}

	return nil
}

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

	//TODO: Переименовать этот метод
	err = createDatabaseIfNotExists(dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalln("Failed to create database:", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

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
