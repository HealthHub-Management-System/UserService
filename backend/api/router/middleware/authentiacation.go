package middleware

import (
	"backend/api/resource/users"
	"fmt"
	"net/http"

	"github.com/wader/gormstore/v2"
)

func checkAPIkey(r *http.Request, apiKey string) bool {
	return r.Header.Get("Authorization") == fmt.Sprintf("Bearer %s", apiKey)
}

func AdminOnly(store *gormstore.Store, apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !checkAPIkey(r, apiKey) {
				session, err := store.Get(r, "session")
				if value, ok := session.Values["role"].(string); !(ok && err == nil && value == users.Admin.ToString()) {
					http.Error(w, "Admin needed!", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func LoggedOnly(store *gormstore.Store, apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !checkAPIkey(r, apiKey) {
				session, err := store.Get(r, "session")
				if _, ok := session.Values["email"].(string); !ok || err != nil {
					http.Error(w, "You must log in!", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
