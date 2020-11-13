package routes

import (
	"encoding/json"
	"github.com/finitum/aurum/internal/aurum"
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

	err := rs.au.SignUp(r.Context(), u)
	if err != nil {
		_ = AutomaticRenderError(w, err)
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

	tp, err := rs.au.Login(r.Context(), u)
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

	err := rs.au.RefreshToken(&tp)
	if err == aurum.ErrInvalidInput {
		_ = AutomaticRenderError(w, err)
		return
	}

	tp.RefreshToken = ""

	_ = json.NewEncoder(w).Encode(&tp)
}

func (rs Routes) GetMe(w http.ResponseWriter, r *http.Request) {
	token := TokenFromContext(r.Context())

	user, err := rs.au.GetUser(r.Context(), token)
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(user)
}

func (rs Routes) SetUser(w http.ResponseWriter, r *http.Request) {
	token := TokenFromContext(r.Context())

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		_ = RenderError(w, err, InvalidRequest)
		return
	}

	user, err := rs.au.UpdateUser(r.Context(), token, u)
	if err != nil {
		_ = AutomaticRenderError(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(&user)
}
