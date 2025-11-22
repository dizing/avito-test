package functional

import (
	"avito-test/internal/handler"
	"avito-test/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthShouldRegisterFirstAndLoginLater(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	tests := []struct {
		name     string
		username string
		password string
	}{
		{"should register new user", "user", "password"},
		{"should login existing user", "user", "password"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testAuthorizeAndGetInfo(t, env, test.username, test.password)
		})
	}
}

// TODO: use something like ginkgo to better describe tests logic
func TestAuthNewUserShouldHaveEmptyInfoWith5000Balance(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	expected_info := handler.InfoResponse{
		Coins:       5000,
		Inventory:   []handler.Inventory{},
		CoinHistory: handler.CoinHistory{},
	}

	info := testAuthorizeAndGetInfo(t, env, "user", "password")

	assert.Equal(t, expected_info, info)
}

func testAuthorizeAndGetInfo(t *testing.T, env *test.Environment, username string, password string) handler.InfoResponse {
	code, token := env.ShopClient.Auth(username, password)
	assert.Equal(t, 200, code)

	env.ShopClient.SetAuthToken(token.Token)

	code, info := env.ShopClient.GetInfo()
	assert.Equal(t, 200, code)

	assert.True(t, info.Coins > 0)

	return info
}

func TestInvalidRequestIfPasswordIncorrect(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	username := "user"
	password := "pass"
	incorrectPassword := "incorrectPass"

	testAuthorizeAndGetInfo(t, env, username, password)

	code, _ := env.ShopClient.Auth(username, incorrectPassword)
	assert.Equal(t, 400, code)
}
