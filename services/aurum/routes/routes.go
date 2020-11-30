package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	log "github.com/sirupsen/logrus"
)

type Routes struct {
	au  aurum.Aurum
	cfg *config.Config
}

func NewRoutes(au aurum.Aurum, cfg *config.Config) Routes {
	return Routes{au, cfg}
}

type ErrorCode int

const (
	ServerError ErrorCode = iota
	InvalidRequest
	Duplicate
	WeakPassword
	Unauthorized
	NotFound
)

type ErrorResponse struct {
	Message string
	Code    ErrorCode
}

func AutomaticRenderError(w http.ResponseWriter, err error) error {
	code := ServerError
	switch err {
	case store.ErrExists:
		code = Duplicate
	case store.ErrNotExists:
		code = NotFound
	case aurum.ErrInvalidInput:
		code = InvalidRequest
	case aurum.ErrWeakPassword:
		code = WeakPassword
	case aurum.ErrUnauthorized:
		code = Unauthorized
	}

	return RenderError(w, err, code)
}

func RenderError(w http.ResponseWriter, err error, code ErrorCode) error {
	switch code {
	case NotFound:
		w.WriteHeader(http.StatusNotFound)
	case Duplicate:
		w.WriteHeader(http.StatusConflict)
	case Unauthorized:
		w.WriteHeader(http.StatusUnauthorized)
	case InvalidRequest, WeakPassword:
		w.WriteHeader(http.StatusBadRequest)
	case ServerError:
		fallthrough
	default:
		log.Errorf("Internal Server Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(&ErrorResponse{
		Message: err.Error(),
		Code:    code,
	}); err != nil {
		return err
	}

	return nil
}

func (rs Routes) PublicKey(w http.ResponseWriter, r *http.Request) {
	pem, err := rs.cfg.PublicKey.ToPem()
	if err != nil {
		_ = RenderError(w, err, ServerError)
	}

	_ = json.NewEncoder(w).Encode(&models.PublicKeyResponse{
		PublicKey: pem,
	})
}

type aurumContextKey string

const (
	contextKeyToken aurumContextKey = "aurum web context key token"
)

// TokenExtractionMiddleware extracts the Authorization token from the http request and stores it in the request context
// you can access this token using the TokenFromContext helper
func (rs Routes) TokenExtractionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyToken, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TokenFromContext(ctx context.Context) string {
	val, ok := ctx.Value(contextKeyToken).(string)
	if !ok {
		return ""
	}
	return val
}
