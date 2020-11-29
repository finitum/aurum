package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"net/http"
)

func AddGroup(host string, tp *jwt.TokenPair, group *models.Group) error {
	body, err := json.Marshal(group)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, host+"/group", bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

func RemoveGroup(host string, tp *jwt.TokenPair, group string) error {
	req, err := http.NewRequest(http.MethodDelete, host+"/group/"+group, nil)
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

const groupUserFmtUrl = "%s/group/%s/%s"

func GetAccess(host string, group, user string) (models.AccessStatus, error) {
	if group == "" {
		group = aurum.AurumName
	}

	url := fmt.Sprintf(groupUserFmtUrl, host, group, user)

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
	url := fmt.Sprintf(groupUserFmtUrl, host, access.GroupName, access.Username)

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

func AddUserToGroup(host string, tp *jwt.TokenPair, user, group string) error {
	url := fmt.Sprintf(groupUserFmtUrl, host, group, user)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

func RemoveUserFromGroup(host string, tp *jwt.TokenPair, user, group string) error {
	url := fmt.Sprintf(groupUserFmtUrl, host, group, user)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = authenticatedRequest(req, tp)
	return err
}

func GetGroupsForUser(host string, tp *jwt.TokenPair, user string) ([]models.GroupWithRole, error) {
	url := fmt.Sprintf("%s/user/%s/groups", host, user)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := authenticatedRequest(req, tp)
	if err != nil {
		return nil, err
	}

	var groups []models.GroupWithRole
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, err
	}

	return groups, nil
}
