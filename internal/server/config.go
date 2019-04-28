package server

import (
	"fmt"
	"time"
)

type (
	Config struct {
		DumpHTTP      bool
		Port          string
		QueryCacheTTL time.Duration
	}

	ConfigOption func(*Config)
)

func (c *Config) ApplyOptions(options ...ConfigOption) {
	for _, option := range options {
		option(c)
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("dump-http: %t, port: %s, querycachettl: %s",
		c.DumpHTTP, c.Port, c.QueryCacheTTL,
	)
}

func DumpHTTP(dump bool) ConfigOption {
	return func(s *Config) {
		s.DumpHTTP = dump
	}
}

func Port(port string) ConfigOption {
	return func(s *Config) {
		s.Port = port
	}
}

func QueryCacheTTL(ttl time.Duration) ConfigOption {
	return func(s *Config) {
		s.QueryCacheTTL = ttl
	}
}
