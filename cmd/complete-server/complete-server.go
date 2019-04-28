package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/reinerRubin/complete-apu-apustaja/internal/server"
)

func main() {
	if err := runApp(); err != nil {
		log.Fatalf("app is terminated with error: %s", err)
	}
}

func runApp() error {
	serverOptions, err := getServerOptionsFromEnv()
	if err != nil {
		return err
	}

	server, err := server.New(serverOptions...)
	if err != nil {
		return err
	}

	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt, syscall.SIGTERM)

	serverErrChan := server.Start()

	select {
	case <-interrupted:
		if err := server.Stop(); err != nil {
			return err
		}
	case err := <-serverErrChan:
		return err
	}

	log.Println("see you space cowboy!")
	return nil
}

func getServerOptionsFromEnv() ([]server.ConfigOption, error) {
	options := make([]server.ConfigOption, 0)

	if port := os.Getenv("PORT"); port != "" {
		_, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("cant parse port (%s): %s", port, err)
		}
		options = append(options, server.Port(port))
	}

	// TODO: migrate to fancy lib
	if debug := os.Getenv("DUMP_HTTP"); debug != "" {
		off := debug == "0" || debug == "false" || debug == "nopls"
		options = append(options, server.DumpHTTP(!off))
	}

	if ttl := os.Getenv("QUERY_CACHE_TTL_SECONDS"); ttl != "" {
		ttlSeconds, err := strconv.Atoi(ttl)
		if err != nil {
			return nil, fmt.Errorf("cant parse port (%s): %s", ttl, err)
		}
		if ttlSeconds <= 0 {
			return nil, fmt.Errorf("query ttl is invalid: %d", ttlSeconds)
		}

		options = append(options, server.QueryCacheTTL(time.Duration(ttlSeconds)*time.Second))
	}

	return options, nil
}
