package main

import (
	"context"
	"log"
	"net/http"

	jsongenerics "github.com/zspekt/rssAggregator/internal/jsonGenerics"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := authWithApiKey(r)
		if err != nil {
			log.Printf("Error authenticating with api key -> %v\n", err)
			jsongenerics.RespondWithError(w, 400, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func secMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}
