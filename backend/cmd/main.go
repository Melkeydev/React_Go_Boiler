package main

import (
  "os"
  "log"
  "fmt"
  "time"
  "flag"
  "net/http"
  "backend/types"
)

type application struct {
  config types.Config
  logger *log.Logger
}

func main() {
  var cfg types.Config
  var port = 4000
  logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

  fmt.Println("Loading server...")

  flag.IntVar(&cfg.Port, "port", port, "server for port to listen")
  flag.StringVar(&cfg.Env, "env", "development", "app environment")
  flag.Parse()
  
  app := &application {
    config: cfg,
    logger:  logger,
  } 

  // Declare Server config
  server := http.Server{
    Addr: fmt.Sprintf(":%d", cfg.Port),
    Handler: app.routes(),
    IdleTimeout: time.Minute,
    ReadTimeout: 10*time.Second,
    WriteTimeout: 30*time.Second,
  }
  
  // Run the server
  fmt.Printf("Server running on port %d", port)
  err := server.ListenAndServe()
  if err != nil {
    log.Println(err)
  }

}
