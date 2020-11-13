package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"net/http"
)

func AddApplication(host string, tp *jwt.TokenPair, app *models.Application) error {
	body, err := json.Marshal(app)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, host+"/application", bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

func RemoveApplication(host string, tp *jwt.TokenPair, app string) error {
	req, err := http.NewRequest(http.MethodDelete, host+"/application/"+app, nil)
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

const applicationUserFmtUrl = "%s/application/%s/%s"

func GetAccess(host string, app, user string) (models.AccessStatus, error) {
	url := fmt.Sprintf(applicationUserFmtUrl, host, app, user)

	resp, err := http.Get(url)
	if err != nil {
		return models.AccessStatus{}, err
	}

	var access models.AccessStatus
	if err := json.NewDecoder(resp.Body).Decode(&access); err != nil {
		return models.AccessStatus{}, err
	}

	return access, nil
}

func SetAccess(host string, tp *jwt.TokenPair, access models.AccessStatus) error {
	url := fmt.Sprintf(applicationUserFmtUrl, host, access.ApplicationName, access.Username)

	body, err := json.Marshal(&access)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

func AddUserToApplication(host string, tp *jwt.TokenPair, user, app string) error {
	url := fmt.Sprintf(applicationUserFmtUrl, host, app, user)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

func RemoveUserFromApplication(host string, tp *jwt.TokenPair, user, app string) error {
	url := fmt.Sprintf(applicationUserFmtUrl, host, app, user)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}
