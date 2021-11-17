package main

import (
	"backend/jsonlog"
	"backend/mailer"
	"backend/models"
	"backend/types"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

type application struct {
	config types.Config
	logger *jsonlog.Logger
	models models.Models
	mailer mailer.Mailer
}

func main() {
	var cfg types.Config
	var port = 4000
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	logger.PrintInfo("Loading server...", nil)

	flag.IntVar(&cfg.Port, "port", port, "server for port to listen")
	flag.StringVar(&cfg.Env, "env", "development", "app environment")
	// TODO: Add to note to the readme
	// CHANGE DSN to your database setting
	flag.StringVar(&cfg.Db.Dsn, "dsn", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable", "Database connection string")
	flag.StringVar(&cfg.Jwt.Secret, "jwt-secret", "default-secret", "secret-key")

	// create flags
	flag.StringVar(&cfg.SMTP.Host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", 587, "SMTP Port")
	// I need to actually put in my credentials
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", "put-your-credentials", "SMTP host")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", "put-your-credentials", "SMTP host")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", "Thundercock <no-reply@thuder.cock.net>", "SMTP host")

	flag.Parse()

	db, err := connectDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
		mailer: mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
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
	logger.PrintInfo("Server running on port", nil)
	err = server.ListenAndServe()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
