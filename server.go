package main

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	srv     *http.Server
	logger  *log.Logger
	closeCh chan<- *Report
}

func NewServer(logger *log.Logger, closeCh chan<- *Report) (*Server, error) {
	mux := http.NewServeMux()
	h, err := newHandler(closeCh)
	if err != nil {
		return nil, err
	}
	mux.Handle("/report", h)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return &Server{
		srv:     srv,
		logger:  logger,
		closeCh: closeCh,
	}, nil
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		s.logger.Println("goroutine started")
		<-ctx.Done()
		// TODO: ctx
		s.Shutdown(context.Background())
	}()

	// handle signals
	sigChan := createSigCh()
	go func() {
		<-sigChan
		s.Shutdown(context.Background())
	}()

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Printf("server failed: %s", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("graceful stopped")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	s.closeCh <- nil
	return nil
}
