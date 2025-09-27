package functional

import (
	"avito-test/test"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuyExistedItem(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	env.EnsureAuthorized()

	default_items := []struct {
		name           string
		expected_price int
	}{
		{"t-shirt", 80},
		{"cup", 20},
		{"book", 50},
		{"pen", 10},
		{"powerbank", 200},
		{"hoody", 300},
		{"umbrella", 200},
		{"socks", 10},
		{"wallet", 50},
		{"pink-hoody", 500},
	}

	for _, item := range default_items {
		t.Run(fmt.Sprintf("Should buy existing item %s with price %d", item.name, item.expected_price), func(t *testing.T) {
			expected_balance := env.GetBalance() - item.expected_price

			env.ShopClient.BuyItem(item.name)

			assert.Equal(t, expected_balance, env.GetBalance())
		})
	}
}

func TestBuyUnexistedItemShouldReturnInvalidRequest(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()
	env.EnsureAuthorized()

	code := env.ShopClient.BuyItem("unexists")
	assert.Equal(t, 400, code)
}

func TestNotEnoughMoneyShouldReturnInvalidRequest(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()
	env.EnsureAuthorized()

	code := env.ShopClient.BuyItem("hoody")
	assert.Equal(t, 200, code)

	// TODO: change balance manually through database injection
	for balance := env.GetBalance(); balance > 200; balance = env.GetBalance() {
		code = env.ShopClient.BuyItem("umbrella")
		assert.Equal(t, 200, code)
	}

	code = env.ShopClient.BuyItem("hoody")
	assert.Equal(t, 400, code)
}

func TestNotAuthorizedBuyShouldReturn401(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	code := env.ShopClient.BuyItem("hoody")
	assert.Equal(t, 401, code)
}
