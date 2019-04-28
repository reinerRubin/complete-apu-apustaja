package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/reinerRubin/complete-apu-apustaja/internal/cache"
)

type Server struct {
	httpServer *http.Server

	errChan chan error

	stopOnce       sync.Once
	stopChannel    chan struct{}
	stoppedChannel chan struct{}

	cache *cache.InMemoryCache
}

func New(options ...ConfigOption) (*Server, error) {
	serverConfig := &Config{
		Port:          "7377",
		DumpHTTP:      true,
		QueryCacheTTL: 10 * time.Second,
	}
	serverConfig.ApplyOptions(options...)
	log.Printf("config: %s", serverConfig)

	cache := cache.NewInMemoryCache(10 * time.Second)
	mux, err := NewServerMux(&serverHandlerContext{cache: cache}, serverConfig)
	if err != nil {
		return nil, err
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverConfig.Port),
		Handler: mux,
	}
	return &Server{
		httpServer: httpServer,

		stopChannel:    make(chan struct{}),
		stoppedChannel: make(chan struct{}),
		errChan:        make(chan error),

		cache: cache,
	}, nil
}

func (s *Server) Start() <-chan error {
	s.cache.Start()
	go s.run()
	return s.errChan
}

func (s *Server) Stop() error {
	s.stopOnce.Do(func() {
		close(s.stopChannel)
	})
	<-s.stoppedChannel

	return nil
}

func (s *Server) run() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("server has stopped: %s", err)
			s.errChan <- err
		}
	}()

	<-s.stopChannel
	s.stopRoutine()
	close(s.stoppedChannel)

	return nil
}

func (s *Server) stopRoutine() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Printf("error on shutdown: %s", err)
	}

	err = s.cache.Stop()
	if err != nil {
		log.Printf("error on cache stop: %s", err)
	}
}
