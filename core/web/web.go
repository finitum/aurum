package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/jwt"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Endpoints struct {
	conn   db.UserRepository
	config *config.Config
}

type contextKey string

func (c contextKey) String() string {
	return "aurum web context key " + string(c)
}

var (
	contextKeyUser   = contextKey("user")
	contextKeyClaims = contextKey("claims")
)

// Set access control headers on all requests
func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Authenticates a HTTP request by verifying the JWT Token
func (e *Endpoints) authenticationMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := jwt.VerifyJWT(token, e.config)
		if err != nil {
			http.Error(w, "Invalid JWT Token", http.StatusUnauthorized)
			return
		}

		// Refresh tokens are not allowed to be used as authentication
		if claims.Refresh {
			http.Error(w, "Please provide Login Token", http.StatusBadRequest)
			return
		}

		// Get user to check if blocked
		u, err := e.conn.GetUserByName(claims.Username)
		if err != nil {
			http.Error(w, "Error retrieving user from database", http.StatusInternalServerError)
			return
		}

		// If blocked deny
		if u.Blocked {
			http.Error(w, "User blocked", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyClaims, claims)
		ctx = context.WithValue(ctx, contextKeyUser, &u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

// Starts the REST web API
func StartServer(config *config.Config, db db.UserRepository) {

	endpoints := Endpoints{
		conn:   db,
		config: config,
	}

	// Router
	router := mux.NewRouter()
	router.Use(accessControlMiddleware)

	unauthenticatedRouter := router.PathPrefix(config.Path).Subrouter()

	// *WARNING* Unauthenticated routes
	unauthenticatedRouter.HandleFunc("/signup", endpoints.signup).Methods(http.MethodPost, http.MethodOptions)
	unauthenticatedRouter.HandleFunc("/login", endpoints.login).Methods(http.MethodPost, http.MethodOptions)
	unauthenticatedRouter.HandleFunc("/refresh", endpoints.refresh).Methods(http.MethodPost, http.MethodOptions)

	// Authenticated routes (Login/ Token required)
	authenticatedRouter := router.PathPrefix(config.Path).Subrouter()
	authenticatedRouter.Use(endpoints.authenticationMiddleware)

	authenticatedRouter.HandleFunc("/changepassword", endpoints.changePassword).Methods(http.MethodPost, http.MethodOptions)
	authenticatedRouter.HandleFunc("/me", endpoints.getMe).Methods(http.MethodGet, http.MethodOptions)
	authenticatedRouter.HandleFunc("/users", endpoints.getUsers).Methods(http.MethodPost, http.MethodOptions)

	// Create the server
	srv := &http.Server{
		Handler: router,
		Addr:    config.WebAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info("Starting up web server ...")
	log.Fatal(srv.ListenAndServe())
}
