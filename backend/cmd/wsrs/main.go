package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/alexandrecpedro/ama-room/backend/internal/api"
	"github.com/alexandrecpedro/ama-room/backend/internal/store/pgstore"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// (1) Load .env variables
	if err := godotenv.Load(); err != nil {
		// panic(err)
		panic(fmt.Sprintf("Error while loading .env: %v", err))
	}

	// (2) DB connection
	// context
	ctx := context.Background()

	// connection pool (pgx manage connections)
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WSRS_DATABASE_USER"),
		os.Getenv("WSRS_DATABASE_PASSWORD"),
		os.Getenv("WSRS_DATABASE_HOST"),
		os.Getenv("WSRS_DATABASE_PORT"),
		os.Getenv("WSRS_DATABASE_NAME"),
	))
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to db: %v", err))
	}

	// (3) Before function return, execute the code described with defer
	// similar to make a clean-up
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("Error while ping connection: %v", err))
	}

	// pgstore.New(DB_CONNECTION) => method created by sqlc
	handler := api.NewHandler(pgstore.New(pool))

	// (4) Async start http server
	// http server is blocking - runs infinitely until the server runs into error
	go func() {
		port := os.Getenv("SERVER_PORT")
		if port == "" {
			port = "8080"
		}
		serverPort := fmt.Sprintf(":%s", port)
		fmt.Println("Server running on port", port)
		err := http.ListenAndServe(serverPort, handler)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(fmt.Sprintf("Error while starting server: %v", err))
			}
		}

	}()

	// (5) Quit when receives interrupt signal from Operational System
	// os.Signal must be buffered
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
