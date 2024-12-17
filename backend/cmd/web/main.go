package main

import (
	"backend/repository"
	"backend/service"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/redis/go-redis/v9"
)

var dbConn *sql.DB
var redisClient *redis.Client

// automatically called before main()
func init() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Channel for error reporting
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		var err error
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Println("DATABASE_URL not set")
			databaseURL = "user:password@tcp(127.0.0.1:3306)/reservations"
		}
		dbConn, err = initDbConnection(databaseURL, 3, time.Second*1)
		if err != nil {
			errChan <- fmt.Errorf("MySQL connection error: %v", err)
			return
		}
		fmt.Println("MySQL connected successfully")
	}()

	go func() {
		defer wg.Done()
		var err error
		addr := os.Getenv("REDIS_URL")
		redisClient, err = initRedisConnection(addr)
		if err != nil {
			errChan <- fmt.Errorf("redis connection error: %v", err)
			return
		}
		fmt.Println("Redis connected successfully")
	}()

	wg.Wait()
	close(errChan)

	// Handle any errors
	for err := range errChan {
		log.Printf("Error: %v\n", err)
	}

	if dbConn == nil {
		log.Fatal("MySQL connection is not initialized")
	}
	if redisClient == nil {
		log.Fatal("Redis client is not initialized")
	}
}

func initDbConnection(dsn string, retries int, delay time.Duration) (*sql.DB, error) {
	var conn *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		conn, err = sql.Open("mysql", dsn)
		if err == nil {
			if pingErr := conn.Ping(); pingErr == nil {
				break
			}
		}
		slog.Info("MySQL not ready, retrying...", "delay", delay)
		time.Sleep(delay)
	}

	slog.Info("Connected to MySQL database\n", "addr", dsn)

	// var rows *sql.Rows
	// rows, err = conn.Query("show databases;")

	// if err != nil {
	// 	log.Fatalf("Failed to execute query: %v", err)
	// }
	// defer rows.Close()

	// var databaseName string
	// fmt.Println("Databases:")
	// for rows.Next() {
	// 	rows.Scan(&databaseName)
	// 	fmt.Println(databaseName)
	// }

	return conn, nil
}

func initRedisConnection(addr string) (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     addr, // Redis server address
			Password: "",   // No password for default setup
			DB:       0,    // Use default DB
		})

	// Test connection to Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	slog.Info("Connected to Redis")

	return redisClient, nil
}

func main() {
	defer cleanup()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
  slog.SetDefault(logger)

  app := &App {
    Database: dbConn,
		Service: service.New(repository.New(dbConn)),
		Router: gin.Default(),
    Redis: redisClient,
    Port: 8080,
		Logger: logger,
  }

	app.Router.GET("/api/reservations", app.HandleGetReservations)
	app.Router.POST("/api/reservations", app.HandlePostReservation)
	app.Router.GET("/api/reservations/:id", app.HandleGetReservationByID)
	// http.HandleFunc("/api/reservations", app.HandleGetReservations)
	// http.HandleFunc("/api/reservations", app.HandlePostReservation)

  app.RunServer()
}

func cleanup() {
	if dbConn != nil {
		err := dbConn.Close()
		if err != nil {
			return
		}
		slog.Info("MySQL connection closed")
	}

	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			return
		}
		slog.Info("Redis connection closed")
	}
}
