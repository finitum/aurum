package api

import (
	"encoding/json"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func GetPublicKey(host string) (*models.PublicKeyResponse, error) {
	resp, err := http.Get(host + "/pk")
	if err != nil {
		return nil, errors.Wrap(err, "error getting public key")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.Wrapf(err, "error connecting (%v), (%v)", resp.StatusCode, body)
	}

	var pk models.PublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&pk); err != nil {
		return nil, err
	}

	return &pk, nil
}

func authenticatedRequest(req *http.Request, tp *jwt.TokenPair) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+tp.LoginToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		if err = Refresh(req.URL.Scheme+"://"+req.URL.Host, tp); err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+tp.LoginToken)
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		return nil, err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return resp, errors.Errorf("unexpected status code (%d)", resp.StatusCode)
	}

	return resp, nil
}
