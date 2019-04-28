package server

import (
	"net/http"
)

type (
	middleware func(http.HandlerFunc) http.HandlerFunc
)

func combineMiddlewares(wrappers ...middleware) middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		accHandler := h
		for i := len(wrappers) - 1; i >= 0; i-- {
			wrapper := wrappers[i]
			accHandler = wrapper(accHandler)
		}

		return accHandler
	}
}
