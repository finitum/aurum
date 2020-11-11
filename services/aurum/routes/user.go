package routes

import (
	"encoding/json"
	"errors"
	"github.com/finitum/aurum/pkg/aurum"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"net/http"
)

func (rs Routes) SignUp(w http.ResponseWriter, r *http.Request) {
	var u models.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	err := aurum.SignUp(r.Context(), rs.store, u)
	if err == aurum.ErrWeakPassword {
		_ = RenderError(w, err, WeakPassword)
		return
	} else if err == aurum.ErrInvalidInput {
		_ = RenderError(w, err, InvalidRequest)
		return
	} else if err != nil {
		_ = RenderError(w, err, ServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rs Routes) Login(w http.ResponseWriter, r *http.Request) {
	var u models.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	tp, err := aurum.Login(r.Context(), rs.store, u, rs.cfg.SecretKey)
	if err != nil {
		_ = RenderError(w, err, Unauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(&tp)
}


func (rs Routes) Refresh(w http.ResponseWriter, r *http.Request) {
	var tp jwt.TokenPair

	if err := json.NewDecoder(r.Body).Decode(&tp); err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	err := aurum.RefreshToken(&tp, rs.cfg.PublicKey, rs.cfg.SecretKey)
	if err == aurum.ErrInvalidInput {
		_ = RenderError(w, err, InvalidRequest)
		return
	} else if err != nil {
		_ = RenderError(w, err, Unauthorized)
		return
	}

	tp.RefreshToken = ""

	_ = json.NewEncoder(w).Encode(&tp)
}

func (rs Routes) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		_ = RenderError(w, errors.New("token claims got lost"), ServerError)
		return
	}

	user, err := aurum.GetUser(r.Context(), rs.store, claims.Username)
	if err != nil {
		_ = RenderError(w, err, ServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(user)
}

func (rs Routes) SetUser(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		_ = RenderError(w, errors.New("token claims got lost"), ServerError)
		return
	}

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	u.Username = claims.Username

	err = aurum.UpdateUser(r.Context(), rs.store, u)
	if err != nil {
		_ = RenderError(w, err, ServerError)
		return
	}
}
