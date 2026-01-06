package httpkit

import (
	"fmt"
	"net"
	"net/http"
	"strings"
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

// CacheMiddleware adds CacheControl headers with max-age set to t
func CacheMiddleware(t time.Duration, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderCacheControl, fmt.Sprintf("max-age=%d", int(t.Seconds())))
		next.ServeHTTP(w, r)
	}
}

// RealIPMiddleware replaces r.RemoteAddr with IP from X-Fowarded-For header if it comes from address contained in tn
func RealIPMiddleware(tn *net.IPNet, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		currIP, _, _ := strings.Cut(r.RemoteAddr, ":")

		if xff := r.Header.Get(HeaderXForwardedFor); xff != "" && tn != nil && tn.Contains(net.ParseIP(currIP)) {
			ip, _, _ = strings.Cut(xff, ",")
		}
		if ip == "" || net.ParseIP(ip) == nil {
			ip, _, _ = strings.Cut(r.RemoteAddr, ":")
		}
		r.RemoteAddr = ip
		next.ServeHTTP(w, r)
	})
}
