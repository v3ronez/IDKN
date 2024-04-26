package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

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

func (app *application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(1, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimitPerClient(next http.Handler) http.Handler {
	type client struct {
		limit    *rate.Limiter
		lastSeen time.Time
	}

	var (
		mut     sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mut.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mut.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		mut.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limit: rate.NewLimiter(2, 4)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limit.Allow() {
			mut.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}
		mut.Unlock()
		next.ServeHTTP(w, r)
	})
}
