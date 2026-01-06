package httpkit

import (
	"fmt"
	"net/http"
	"time"
)

const (
	HeaderAccept         = "Accept"
	HeaderAcceptEncoding = "Accept-Encoding"
	HeaderAcceptLanguage = "Accept-Language"
	HeaderCacheControl   = "Cache-Control"
	HeaderContentType    = "Content-Type"
	HeaderReferer        = "Referer"
	HeaderUserAgent      = "User-Agent"
	HeaderXForwardedFor  = "X-Forwarded-For"

	MIMEApplicationJSON = "application/json"
	MIMETextHTML        = "text/html"
)

// ReadyHandler always return StatusOK(200) as response
func ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func CacheMiddleware(t time.Duration, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderCacheControl, fmt.Sprintf("max-age=%d", int(t.Seconds())))
		next.ServeHTTP(w, r)
	}
}
