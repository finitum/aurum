package routes

import (
	"context"
	"encoding/json"
	"github.com/finitum/aurum/pkg/aurum"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"net/http"
	"strings"
)

type Routes struct {
	store store.AurumStore
	cfg   *config.Config
}

func NewRoutes(s store.AurumStore, cfg *config.Config) Routes {
	return Routes{s, cfg}
}

type ErrorCode int

const (
	ServerError ErrorCode = iota
	InvalidRequest
	WeakPassword
	Unauthorized
)

type ErrorResponse struct {
	Message string
	Code    ErrorCode
}

func RenderError(w http.ResponseWriter, err error, code ErrorCode) error {
	switch code {
	case Unauthorized:
		w.WriteHeader(http.StatusUnauthorized)
	case InvalidRequest, WeakPassword:
		w.WriteHeader(http.StatusBadRequest)
	case ServerError:
		fallthrough
	default:
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

const (
	contextKeyClaims = "aurum web context key claims"
	contextKeyRoles  = "aurum web context key roles"
)

// RoleMiddleware depends on TokenVerificationMiddleware. Without it, it will always return unauthorized.
// Rolemiddleware only allows users in with a role that's more powerful than the roll passed in, for the given application.
func (rs Routes) RoleMiddleware(role models.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			claims := ClaimsFromContext(ctx)
			if claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			resp, err := aurum.Access(ctx, rs.store, claims.Username, aurum.Aurum)
			if err != nil || !resp.AllowedAccess || resp.Role < role {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, contextKeyRoles, resp.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (rs Routes) TokenVerificationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := jwt.VerifyJWT(token, rs.cfg.PublicKey)
		if err != nil {
			http.Error(w, "Invalid JWT Token", http.StatusUnauthorized)
			return
		}

		// Refresh tokens are not allowed to be used as authentication
		if claims.Refresh {
			http.Error(w, "Please provide Login Token", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyClaims, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFromContext(ctx context.Context) *jwt.Claims {
	val, ok := ctx.Value(contextKeyClaims).(*jwt.Claims)
	if !ok {
		return nil
	}
	return val
}

func RoleFromContext(ctx context.Context) models.Role {
	val, ok := ctx.Value(contextKeyRoles).(models.Role)
	if !ok {
		return 0
	}
	return val
}
