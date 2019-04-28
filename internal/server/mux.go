package server

import (
	"net/http"
	"time"

	"github.com/reinerRubin/complete-apu-apustaja/internal/cache"
	"github.com/reinerRubin/complete-apu-apustaja/internal/completer"
	"github.com/reinerRubin/complete-apu-apustaja/internal/handler"
)

type serverHandlerContext struct {
	cache cache.Cache
}

// TODO: set handlers via options pattern
// think about serverHandlerContext & ServerConfig
func NewServerMux(ctx *serverHandlerContext, config *Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	{ // complete handler
		avCompleter, _ := completer.NewPlacesaviasales()
		completer := completer.NewCacheableCompleter(
			avCompleter,
			ctx.cache,
			config.QueryCacheTTL,
		)

		completerHandler := combineMiddlewares(
			withTimeout(3*time.Second),
			withLogging(config.DumpHTTP))(handler.NewCompleterHandler(completer).Handle)
		mux.HandleFunc("/complete", completerHandler)
	}

	return mux, nil
}
