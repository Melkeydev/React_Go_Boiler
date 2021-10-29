package main

import (
	"backend/models"
	"backend/types"
	"context"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	config types.Config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg types.Config
	var port = 4000
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	fmt.Println("Loading server...")

	flag.IntVar(&cfg.Port, "port", port, "server for port to listen")
	flag.StringVar(&cfg.Env, "env", "development", "app environment")
	// TODO: Add to note to the readme
	// CHANGE DSN to your database setting
	flag.StringVar(&cfg.Db.Dsn, "dsn", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable", "Database connection string")
	flag.StringVar(&cfg.Jwt.Secret, "jwt-secret", "default-secret", "secret-key")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := connectDB(ctx, cfg)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	// Declare Server config
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Run the server
	fmt.Printf("Server running on port %d", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
