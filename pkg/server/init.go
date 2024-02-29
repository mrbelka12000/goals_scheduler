// Package server is responsible for all server actions.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"goals_scheduler/pkg/config"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

const (
	shutdownTimeout = 3 * time.Second
	timeout         = 15 * time.Second
)

func NewServer(r *mux.Router, cfg config.Config) *Server {

	httpSever := &http.Server{
		Addr:         ":" + cfg.HttpPort,
		Handler:      r,
		ReadTimeout:  timeout,
		IdleTimeout:  timeout,
		WriteTimeout: timeout,
	}

	server := &Server{
		server:          httpSever,
		notify:          make(chan error),
		shutdownTimeout: shutdownTimeout,
	}
	server.start()

	return server
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
