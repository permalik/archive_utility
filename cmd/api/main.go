package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/permalik/utility/db"
)

type config struct {
	env  string
	port int
}

type application struct {
	config config
	ctx    context.Context
	pool   *sql.DB
	logger *log.Logger
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env", err)
	}

	var cfg config
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.port, "port", 5555, "Network port (default 5555)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pool := db.InitDB()
	defer pool.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app := &application{
		config: cfg,
		ctx:    ctx,
		pool:   pool,
		logger: logger,
	}

	if err := db.Ping(ctx); err != nil {
		logger.Fatal(err)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.Router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Printf("starting %s server on port %s\n", cfg.env, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	logger.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server shutdown failed:", err)
	}

	logger.Println("server exited gracefully")
}
