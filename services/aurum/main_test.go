package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/finitum/aurum/clients/go"
	internal "github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/store/dgraph"
	"github.com/stretchr/testify/assert"

	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
)

const url = "http://localhost:8042"

func VerifyLogin(assert *assert.Assertions, client aurum.Client, u models.User) jwt.TokenPair {
	tp, err := client.Login(u.Username, u.Password)
	assert.NoError(err)

	return *tp
}

func VerifySignupLogin(assert *assert.Assertions, client aurum.Client, u models.User) jwt.TokenPair {
	err := client.Register(u.Username, u.Password, u.Email)
	assert.NoError(err)

	return VerifyLogin(assert, client, u)
}

func VerifyGetUser(assert *assert.Assertions, client aurum.Client, tp jwt.TokenPair, expected models.User) {
	user, err := client.GetUserInfo(&tp)
	assert.NoError(err)

	assert.Equal(expected.Username, user.Username)
	assert.Equal(expected.Email, user.Email)
}

func VerifyRefresh(assert *assert.Assertions, client aurum.Client, tp jwt.TokenPair, u models.User, pk ecc.PublicKey) {
	oldClaims, err := jwt.VerifyJWT(tp.LoginToken, pk)
	assert.NoError(err)

	body, err := json.Marshal(tp)
	assert.NoError(err)

	// Wait so that the refresh token definitely should have a higher iat
	time.Sleep(2 * time.Second)

	// Refresh
	resp, err := http.Post(url+"/refresh", "application/json", bytes.NewBuffer(body))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var rtp jwt.TokenPair
	err = json.NewDecoder(resp.Body).Decode(&rtp)
	assert.Empty(rtp.RefreshToken)

	newClaims, err := jwt.VerifyJWT(rtp.LoginToken, pk)
	assert.NoError(err)

	assert.True(oldClaims.IssuedAt < newClaims.IssuedAt)

	tp.LoginToken = rtp.LoginToken
	VerifyGetUser(assert, client, tp, u)
}

func VerifyUpdateUserPasswordEmail(assert *assert.Assertions, client aurum.Client, tp jwt.TokenPair, u models.User) {
	newuser := models.User{
		Username: u.Username,
		Password: "9054fbe0b622c638224d50d20824d2ff6782e308",
		Email:    "yeet42@finitum.dev",
	}

	resp, err := client.UpdateUser(&tp, &newuser)
	assert.NoError(err)

	assert.Equal(u.Username, resp.Username)
	assert.Equal(newuser.Email, resp.Email)
	assert.Empty(resp.Password)

	u.Password = newuser.Password
	u.Email = newuser.Email

	time.Sleep(2 * time.Second)

	VerifyLogin(assert, client, u)
	VerifyGetUser(assert, client, tp, u)
}

func VerifyGetGroupsForUser(assert *assert.Assertions, client aurum.Client, tp jwt.TokenPair, u models.User, expected models.GroupWithRole) {
	groups, err := client.GetGroupsForUser(&tp, u.Username)
	assert.NoError(err)

	assert.Contains(groups, expected)
}

func VerifyAccess(assert *assert.Assertions, client *aurum.RemoteClient, group models.Group, user models.User, role models.Role) {
	access, err := client.GetAccess(group.Name, user.Username)
	assert.NoError(err)

	assert.Equal(group.Name, access.GroupName)
	assert.Equal(user.Username, access.Username)
	assert.True(access.AllowedAccess)
	assert.Equal(role, access.Role)
}

func VerifyNoAccess(assert *assert.Assertions, client *aurum.RemoteClient, group models.Group, user models.User) {
	access, err := client.GetAccess(group.Name, user.Username)
	assert.NoError(err)

	assert.Equal(group.Name, access.GroupName)
	assert.Equal(user.Username, access.Username)
	assert.False(access.AllowedAccess)
	assert.Equal(models.Role(0), access.Role)
}

func MakeUserAdminDirectly(assert *assert.Assertions, username string) {
	ctx := context.Background()
	dg, err := dgraph.New(ctx, "localhost:9080")
	assert.NoError(err)

	err = dg.AddGroupToUser(ctx, username, internal.AurumName, models.RoleAdmin)
	assert.NoError(err)
}

func TestSystemIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	assert := assert.New(t)

	assert.NoError(os.Setenv("NO_KEY_WRITE", "true"))
	assert.NoError(os.Setenv("WEB_ADDRESS", strings.TrimPrefix(url, "http://")))

	// Startup the server
	go main()

	// totally not flaky or something
	// Wait for the server  to start up
	time.Sleep(5 * time.Second)

	client, err := aurum.NewRemoteClient(url)
	assert.NoError(err)

	userOne := models.User{
		Username: "UserOne",
		Email:    "Tester@test.com",
		Password: "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
	}

	userTwo := models.User{
		Username: "UserTwo",
		Email:    "Tester@test.com",
		Password: "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
	}

	// Now run all the endpoint verifications

	resp, err := http.Get(url + "/pk")
	assert.NoError(err)

	var r models.PublicKeyResponse

	err = json.NewDecoder(resp.Body).Decode(&r)
	assert.NoError(err)

	pk, err := ecc.FromPem([]byte(r.PublicKey))
	assert.NoError(err)

	pub := pk.(ecc.PublicKey)

	tpUserOne := VerifySignupLogin(assert, client, userOne)
	tpUserTwo := VerifySignupLogin(assert, client, userTwo)

	time.Sleep(time.Second * 3)

	MakeUserAdminDirectly(assert, userOne.Username)

	VerifyGetUser(assert, client, tpUserOne, userOne)
	VerifyGetUser(assert, client, tpUserTwo, userTwo)

	VerifyRefresh(assert, client, tpUserOne, userOne, pub)
	VerifyRefresh(assert, client, tpUserTwo, userTwo, pub)

	VerifyUpdateUserPasswordEmail(assert, client, tpUserOne, userOne)
	VerifyUpdateUserPasswordEmail(assert, client, tpUserTwo, userTwo)

	// Group tests

	aurumGroup := models.GroupWithRole{
		Group: models.Group{
			Name:              internal.AurumName,
			AllowRegistration: true,
		},
		Role: models.RoleUser,
	}

	VerifyGetGroupsForUser(assert, client, tpUserTwo, userTwo, aurumGroup)
	// Admin
	aurumGroup.Role = models.RoleAdmin
	VerifyGetGroupsForUser(assert, client, tpUserOne, userOne, aurumGroup)

	group := models.Group{
		Name:              "somegroup",
		AllowRegistration: true,
	}

	err = client.AddGroup(&tpUserOne, &group)
	assert.NoError(err)

	time.Sleep(time.Second)

	VerifyAccess(assert, client, group, userOne, models.RoleAdmin)
	VerifyNoAccess(assert, client, group, userTwo)

	err = client.AddUserToGroup(&tpUserOne, userTwo.Username, group.Name)
	assert.NoError(err)

	time.Sleep(time.Second)

	VerifyAccess(assert, client, group, userTwo, models.RoleUser)

	err = client.SetAccess(&tpUserOne, models.AccessStatus{
		GroupName:     group.Name,
		Username:      userTwo.Username,
		AllowedAccess: true,
		Role:          models.RoleAdmin,
	})
	assert.NoError(err)

	time.Sleep(time.Second)

	VerifyAccess(assert, client, group, userTwo, models.RoleAdmin)

	err = client.RemoveUserFromGroup(&tpUserOne, userTwo.Username, group.Name)
	assert.NoError(err)
	time.Sleep(time.Second)
	VerifyNoAccess(assert, client, group, userTwo)

	err = client.RemoveGroup(&tpUserOne, group.Name)
	assert.NoError(err)
	time.Sleep(time.Second)
	VerifyNoAccess(assert, client, group, userOne)
}
