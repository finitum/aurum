package routes

import (
	"encoding/json"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// POST /application (Authenticated)
func (rs Routes) AddApplication(w http.ResponseWriter, r *http.Request) {
	var app models.Application
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	token := TokenFromContext(r.Context())

	if err := rs.au.AddApplication(r.Context(), token, app); err != nil {
		if err == store.ErrExists {
			_ = RenderError(w, err, Duplicate)
			return
		}
		_ = RenderError(w, err, ServerError)
		return
	}

	app.Name = strings.ToLower(app.Name)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&app)
}

// DELETE /application/{name} (Authenticated)
func (rs Routes) RemoveApplication(w http.ResponseWriter, r *http.Request) {
	app := chi.URLParam(r, "name")
	if app == "" {
		_ = RenderError(w, aurum.ErrUnauthorized, Unauthorized)
		return
	}

	token := TokenFromContext(r.Context())

	if err := rs.au.RemoveApplication(r.Context(), token, app); err != nil {
		if err == store.ErrExists {
			_ = RenderError(w, err, Duplicate)
			return
		} else if err == aurum.ErrUnauthorized {
			_ = RenderError(w, err, Unauthorized)
			return
		}

		_ = RenderError(w, err, ServerError)
		return
	}
}

// GET /application/{app}/{user}
func (rs Routes) GetAccess(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	name := chi.URLParam(r, "app")

	if name == "" || user == "" {
		_ = RenderError(w, errors.New("name empty"), InvalidRequest)
		return
	}

	resp, err := rs.au.GetAccess(r.Context(), user, name)
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(&resp)
}

// PUT /application/{app}/{user}
func (rs Routes) SetAccess(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	app := chi.URLParam(r, "app")
	ctx := r.Context()

	if app == "" || user == "" {
		_ = RenderError(w, errors.New("app empty"), InvalidRequest)
		return
	}

	var access models.AccessStatus
	err := json.NewDecoder(r.Body).Decode(&access)
	if err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	if app != access.ApplicationName || user != access.Username {
		_ = RenderError(w, errors.New("body doesn't match path"), InvalidRequest)
		return
	}

	token := TokenFromContext(ctx)

	if access.AllowedAccess {
		err = rs.au.SetAccess(ctx, token, app, user, access.Role)
	} else {
		err = rs.au.RemoveUserFromApplication(ctx, token, user, app)
	}
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}
}

// POST /application/{app}/{user}
func (rs Routes) AddUserToApplication(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	app := chi.URLParam(r, "app")
	ctx := r.Context()

	if app == "" || user == "" {
		_ = RenderError(w, errors.New("app empty"), InvalidRequest)
		return
	}

	token := TokenFromContext(ctx)
	err := rs.au.AddUserToApplication(ctx, token, user, app, models.RoleUser)
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DELETE /application/{app}/{user}
func (rs Routes) RemoveUserFromApplication(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	app := chi.URLParam(r, "app")
	ctx := r.Context()

	if app == "" || user == "" {
		_ = RenderError(w, errors.New("app empty"), InvalidRequest)
		return
	}

	token := TokenFromContext(ctx)
	err := rs.au.RemoveUserFromApplication(ctx, token, app, user)
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}
}
