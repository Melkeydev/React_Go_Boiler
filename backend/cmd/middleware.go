package main

import (
  "net/http"
  "fmt"
  "golang.org/x/time/rate"
  "time"
  "sync"
  "net"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
    next.ServeHTTP(w, r)
  })
}

// route and handle panics better
func (app *application) recoverPanic(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    defer func() {
      if err := recover(); err != nil {
        w.Header().Set("Connection", "close")
        app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
      }
    }()
    next.ServeHTTP(w, r)
  })

} 

func (app *application) rateLimit(next http.Handler) http.Handler{
  type client struct {
    limiter *rate.Limiter
    lastSeen time.Time
  }

  var (
    mu sync.Mutex
    clients = make(map[string]*client)
  )

  go func() {
    for {
      time.Sleep(time.Minute)
      mu.Lock()

      for ip, client := range clients {
        if time.Since(client.lastSeen) > 3*time.Second {
          delete(clients, ip)
        }
      }
      mu.Unlock()
    }
  }() 

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if app.config.Limiter.Enabled {
      // get the ip addy from each request
      ip, _, err := net.SplitHostPort(r.RemoteAddr)
      if err != nil {
        app.serverErrorResponse(w, r, err)
        return
      }

      mu.Lock()

      if _, found := clients[ip]; !found {
        clients[ip] = &client {
          limiter: rate.NewLimiter(rate.Limit(app.config.Limiter.Rps), app.config.Limiter.Burst),
        }
      }

      // Every new ip that gets added to our clients slice gets a time stamp
      clients[ip].lastSeen = time.Now()

      if !clients[ip].limiter.Allow() {
        mu.Unlock()
        app.rateLimitExceededResponse(w, r)
        return
      }

      mu.Unlock()
    }

    next.ServeHTTP(w, r)
  })
}















