package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// POST /group (Authenticated)
func (rs Routes) AddGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	token := TokenFromContext(r.Context())

	if err := rs.au.AddGroup(r.Context(), token, group); err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}

	group.Name = strings.ToLower(group.Name)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&group)
}

// DELETE /group/{group} (Authenticated)
func (rs Routes) RemoveGroup(w http.ResponseWriter, r *http.Request) {
	group := chi.URLParam(r, "group")
	if group == "" {
		_ = RenderError(w, aurum.ErrInvalidInput, InvalidRequest)
		return
	}

	token := TokenFromContext(r.Context())

	if err := rs.au.RemoveGroup(r.Context(), token, group); err != nil {
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

// GET /group/{group}/{user}
func (rs Routes) GetAccess(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	name := chi.URLParam(r, "group")

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

// PUT /group/{group}/{user}
func (rs Routes) SetAccess(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	group := chi.URLParam(r, "group")
	ctx := r.Context()

	if group == "" || user == "" {
		_ = RenderError(w, errors.New("group empty"), InvalidRequest)
		return
	}

	var access models.AccessStatus
	err := json.NewDecoder(r.Body).Decode(&access)
	if err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	if group != access.GroupName || user != access.Username {
		_ = RenderError(w, errors.New("body doesn't match path"), InvalidRequest)
		return
	}

	token := TokenFromContext(ctx)

	if access.AllowedAccess {
		err = rs.au.SetAccess(ctx, token, group, user, access.Role)
	} else {
		err = rs.au.RemoveUserFromGroup(ctx, token, user, group)
	}
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}
}

// POST /group/{group}/{user}
func (rs Routes) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	group := chi.URLParam(r, "group")
	ctx := r.Context()

	if group == "" || user == "" {
		_ = RenderError(w, errors.New("group empty"), InvalidRequest)
		return
	}

	token := TokenFromContext(ctx)
	err := rs.au.AddUserToGroup(ctx, token, user, group, models.RoleUser)
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DELETE /group/{group}/{user}
func (rs Routes) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	group := chi.URLParam(r, "group")
	ctx := r.Context()

	if group == "" || user == "" {
		_ = RenderError(w, errors.New("group empty"), InvalidRequest)
		return
	}

	token := TokenFromContext(ctx)
	err := rs.au.RemoveUserFromGroup(ctx, token, user, group)
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}
}
