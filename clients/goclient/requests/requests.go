package requests

import (
	"encoding/json"
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
		return nil, errors.Wrapf(err, "error connecting (%v)", body)
	}

	var pk models.PublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&pk); err != nil {
		return nil, err
	}

	return &pk, nil
}
