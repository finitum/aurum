package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/jwt"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Endpoints struct {
	conn   db.Connection
	config *config.Config
}

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

// Starts the REST web API
func StartServer(config *config.Config, db db.Connection) {

	endpoints := Endpoints{
		conn:   db,
		config: config,
	}

	// Router
	router := mux.NewRouter()
	router.Use(accessControlMiddleware)

	api := router.PathPrefix(config.Path).Subrouter()

	api.HandleFunc("/signup", endpoints.signup).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/login", endpoints.login).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/refresh", endpoints.refresh).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/changepassword", endpoints.changePassword).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/me", endpoints.getMe).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/requestusers", endpoints.getUsers).Methods(http.MethodPost, http.MethodOptions)

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

// Authenticates a HTTP request by verifying the JWT Token
func (e *Endpoints) authenticateRequest(w http.ResponseWriter, r *http.Request) (*jwt.Claims, error) {

	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
		return nil, errors.New("malformed Authorization header")
	}
	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := jwt.VerifyJWT(token, e.config)
	if err != nil {
		http.Error(w, "Invalid JWT Token", http.StatusUnauthorized)
		return nil, err
	}

	// Refresh tokens are not allowed to be used as authentication
	if claims.Refresh {
		http.Error(w, "Please provide Login Token", http.StatusBadRequest)
		return nil, errors.New("token wasn't a refresh token")
	}

	return claims, nil
}
