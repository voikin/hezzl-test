package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/voikin/hezzl-test/config"
	"github.com/voikin/hezzl-test/internal/controller"
	"github.com/voikin/hezzl-test/internal/repository"
	"github.com/voikin/hezzl-test/internal/service"
)

const configPath = "./config/config.yaml"

func main() {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pg, err := sql.Open("postgres", cfg.Postgres.URL)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer pg.Close()

	if err := pg.Ping(); err != nil {
		log.Fatalf("failed to ping PostgreSQL: %v", err)
	}

	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer func() { _ = nc.Drain() }()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create JetStream context: %v", err)
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Clickhouse.Addr, cfg.Clickhouse.NativePort)},
		Auth: clickhouse.Auth{
			Database: cfg.Clickhouse.DB,
			Username: cfg.Clickhouse.Username,
			Password: cfg.Clickhouse.Password,
		},
	})
	if err != nil {
		log.Fatalf("failed to create ClickHouse connection: %v", err)
	}
	defer conn.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	repos := repository.NewRepositories(pg, conn, js, redisClient)
	services := service.NewServices(repos, js)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go services.EventSaver.Start(ctx)

	ginEngine := gin.Default()
	controller.RegisterRoutes(ginEngine, services)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: ginEngine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ailed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	cancel()
	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
