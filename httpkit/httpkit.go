package httpkit

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	HeaderAccept          = "Accept"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderAcceptLanguage  = "Accept-Language"
	HeaderCacheControl    = "Cache-Control"
	HeaderContentType     = "Content-Type"
	HeaderReferer         = "Referer"
	HeaderUserAgent       = "User-Agent"
	HeaderXForwardedFor   = "X-Forwarded-For"
	HeaderXForwardedProto = "X-Forwarded-Proto"

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

// ProxyHeadersMiddleware sets cliend address and scheme in [http.Request.Context] if request comes from proxy address contained in tn
func ProxyHeadersMiddleware(tn *net.IPNet, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip, scheme string

		currIP, _, _ := strings.Cut(r.RemoteAddr, ":")

		if xff := r.Header.Get(HeaderXForwardedFor); xff != "" && tn != nil && tn.Contains(net.ParseIP(currIP)) {
			ip, _, _ = strings.Cut(xff, ",")
			scheme = r.Header.Get(HeaderXForwardedProto)
		}
		if ip == "" || net.ParseIP(ip) == nil {
			ip, _, _ = strings.Cut(r.RemoteAddr, ":")
		}

		if r.TLS != nil {
			scheme = "https"
		}

		if scheme == "" {
			scheme = "http"
		}

		ctx := context.WithValue(r.Context(), contextRemoteAddr, ip)
		ctx = context.WithValue(ctx, contextScheme, scheme)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
