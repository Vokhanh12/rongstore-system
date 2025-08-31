package auth

import (
	"net/http"
	"strings"
)

func RequireHeaders(rules []HeaderRule) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, rule := range rules {
				if strings.HasPrefix(r.URL.Path, rule.PathPrefix) {
					for _, h := range rule.Required {
						if r.Header.Get(h) == "" {
							http.Error(w, "missing header: "+h, http.StatusUnauthorized)
							return
						}
					}
					break
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
