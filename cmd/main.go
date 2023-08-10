package main

import (
	"be-pokemon-club/internal/database"
	"be-pokemon-club/internal/middleware"
	"be-pokemon-club/internal/redis"
	"be-pokemon-club/internal/routes"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	database.Init()

	redis.Init()

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "svc-pokemon-club",
		AppName:       "service backend pokemon club",
	})

	// Setup logger
	app.Use(middleware.LoggerMiddleware())

	// Setup routes
	routes.Setup(app)

	go func() {
		serverAddr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
		fmt.Printf("Server running on %s\n", serverAddr)

		if err := app.Listen(serverAddr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	//channel for catch signal SIGINT (Ctrl+C) dan SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// wait signal to stop server
	sig := <-sigCh
	fmt.Printf("Received termination signal: %v\n", sig)

	// Create a context with timeout for managing shutdown
	_, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer timeoutCancel()

	// Initiate shutdown of the Fiber app
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	fmt.Println("Server shutdown complete")
}
