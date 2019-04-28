package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)

// TODO log into one line
func withLogging(dumpHTTP bool) middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if !dumpHTTP {
			return next
		}

		return func(w http.ResponseWriter, r *http.Request) {
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			result := rec.Result()
			responseDump, err := httputil.DumpResponse(result, true)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}

			for k, v := range rec.HeaderMap {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			// TODO check if this is approach is ok
			w.Write(rec.Body.Bytes())

			log.Printf("%s -> %s", string(requestDump), string(responseDump))
		}
	}
}
