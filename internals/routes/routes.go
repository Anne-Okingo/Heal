package routes

import (
	"net/http"
	"strings"
)

// Allowed routes
var allowedRoutes = map[string]bool{
	"/":                 true,
	"/privacy":          true,
	"/session":          true,
	"/signup":           true,
	"/loggin":           true,
	"/welcome":          true,
	"/api/get-username": true,
	"/api/logout":       true,
	"/api/gemini":       true,
	"/api/speechify":    true,
	"/get-heard":        true,
}

// RouteChecker is a middleware that checkes allowed routes
func RouteChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			// Static(w,r)
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := allowedRoutes[r.URL.Path]; !ok {
			handlers.NotFoundHandler(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
