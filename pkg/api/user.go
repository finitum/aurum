package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/pkg/errors"
)

func SignUp(host string, user models.User) error {
	userb, err := json.Marshal(&user)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal user")
	}

	resp, err := http.Post(host+"/signup", "group/json", bytes.NewReader(userb))
	if err != nil {
		return errors.Wrap(err, "couldn't post signup request")
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)

		return errors.Errorf("Unexpected status code (%v), (%v)", resp.StatusCode, string(body))
	}

	return nil
}

func Login(host string, user models.User) (*jwt.TokenPair, error) {
	userb, err := json.Marshal(&user)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't marshal user")
	}

	resp, err := http.Post(host+"/login", "group/json", bytes.NewReader(userb))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't post login request")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, errors.Errorf("Unexpected status code (%v): %v", resp.StatusCode, string(body))
	}

	var tp jwt.TokenPair
	if err := json.NewDecoder(resp.Body).Decode(&tp); err != nil {
		return nil, errors.Wrap(err, "couldn't decode json body")
	}

	return &tp, err
}

func Refresh(host string, tp *jwt.TokenPair) error {
	tpb, err := json.Marshal(tp)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal token")
	}

	resp, err := http.Post(host+"/refresh", "group/json", bytes.NewReader(tpb))
	if err != nil {
		return errors.Wrap(err, "couldn't post refresh request")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		return errors.Errorf("Unexpected status code (%v), (%v)", resp.StatusCode, string(body))
	}

	var newtp jwt.TokenPair
	if err := json.NewDecoder(resp.Body).Decode(&newtp); err != nil {
		return errors.Wrap(err, "couldn't decode json body")
	}

	tp.LoginToken = newtp.LoginToken
	return nil
}

func GetUser(host string, tp *jwt.TokenPair) (*models.User, error) {
	req, err := http.NewRequest(http.MethodGet, host+"/user", nil)
	if err != nil {
		return nil, err
	}

	resp, err := authenticatedRequest(req, tp)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(host string, tp *jwt.TokenPair, user *models.User) (ret *models.User, _ error) {
	userb, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling json")
	}

	req, err := http.NewRequest(http.MethodPost, host+"/user", bytes.NewReader(userb))
	if err != nil {
		return nil, errors.Wrap(err, "building update user request")
	}

	resp, err := authenticatedRequest(req, tp)
	if err != nil {
		return nil, errors.Wrap(err, "update user")
	}

	return ret, errors.Wrap(json.NewDecoder(resp.Body).Decode(&ret), "json decoding response")
}
