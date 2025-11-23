package functional

import (
	"avito-test/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendShouldCorrectlyChangeBalance(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()
	env.EnsureTestUserAuthorized()

	another_username, _, another_user_info := env.EnsureHaveAnotherUser()
	expected_receiver_balance := another_user_info.Coins + 50
	expected_user_balance := env.GetTestUserBalance() - 50

	code := env.ShopClient.Send(another_username, 50)
	assert.Equal(t, 200, code)

	_, _, another_user_info = env.EnsureHaveAnotherUser()
	assert.Equal(t, expected_receiver_balance, another_user_info.Coins)
	assert.Equal(t, expected_user_balance, env.GetTestUserBalance())
}

func TestNotAuthorizedSendShouldReturn401(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	another_username, _, _ := env.EnsureHaveAnotherUser()

	code := env.ShopClient.Send(another_username, 20)
	assert.Equal(t, 401, code)
}

func TestSendToNonExistentUserShouldReturn400(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()
	env.EnsureTestUserAuthorized()

	code := env.ShopClient.Send("some random user", 20)
	assert.Equal(t, 400, code)
}
