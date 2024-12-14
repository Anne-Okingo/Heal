package routes

import (
	"net/http"
	"strings"

	"Heal/internals/handlers"
	"Heal/internals/renders"
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
	"/report": true,
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

// RegisterRoutes manages the routes
func RegisterRoutes(mux *http.ServeMux) {
	staticDir := renders.GetProjectRoot("views", "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r)
	})

	mux.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		handlers.DataPrivacyHandler(w, r)
	})

	mux.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		handlers.WelcomeHandler(w, r)
	})

	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		handlers.SessionHandler(w, r)
	})

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignupHandler(w, r)
	})

	mux.HandleFunc("/loggin", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r)
	})

	mux.HandleFunc("/api/get-username", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUsernameHandler(w, r)
	})

	mux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutHandler(w, r)
	})

	mux.HandleFunc("/api/gemini", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProxyGeminiRequest(w, r)
	})

	mux.HandleFunc("/api/speechify", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProxySpeechifyRequest(w, r)
	})

	mux.HandleFunc("/get-heard", func(w http.ResponseWriter, r *http.Request) {
		handlers.ChatHandler(w, r)
	})

	mux.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		handlers.ReportHandler(w, r)
	})
}
