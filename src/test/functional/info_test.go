package functional

import (
	"avito-test/internal/handler"
	"avito-test/test"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotAuthorizedInfoShouldReturn401(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	code, _ := env.ShopClient.GetInfo()
	assert.Equal(t, 401, code)
}

// NOTE: info -> buy/auth/send; env.GetInfo used in another tests through GetBalance and assert ALL other tests correctness
// there we must assert that GetInfo exactly match user in database
func TestGetInfoCorrectlyReturnUserInfo(t *testing.T) {
	env := test.NewEnvironment(t)
	defer env.Close()

	first_user, _, _ := env.EnsureTestUserAuthorized()
	second_user, second_user_token, _ := env.EnsureHaveAnotherUser()

	ctx := context.Background()

	query := `
		INSERT INTO Transactions (username_from, username_to, amount) 
		VALUES ($1, $2, 100);

		UPDATE Users SET balance = balance - 100 WHERE username = $1;
		UPDATE Users SET balance = balance + 100 WHERE username = $2;

		INSERT INTO User_Items (username, item_name, amount) 
		VALUES 
				($1, 'hoody', 1),
				($1, 'powerbank', 1);

		UPDATE Users SET balance = balance - (SELECT price from items where name='hoody') WHERE username = $1;
		UPDATE Users SET balance = balance - (SELECT price from items where name='powerbank') WHERE username = $1;

		INSERT INTO User_Items (username, item_name, amount) 
		VALUES 
				($2, 'umbrella', 1),
				($2, 'wallet', 1);

		UPDATE Users SET balance = balance - (SELECT price from items where name='umbrella') WHERE username = $2;
		UPDATE Users SET balance = balance - (SELECT price from items where name='wallet') WHERE username = $2;
	`
	query = strings.ReplaceAll(query, "$1", fmt.Sprintf("'%s'", first_user))
	query = strings.ReplaceAll(query, "$2", fmt.Sprintf("'%s'", second_user))

	_, err := env.ShopDatabase.Pool.Exec(ctx, query)
	if err != nil {
		log.Fatal("Failed to setup test data: ", err)
	}

	const INITIAL_MONEY = 5000
	const TRANSACTION_MONEY = 100
	itemPrices := make(map[string]int)
	itemPrices["hoody"] = 300
	itemPrices["umbrella"] = 200
	itemPrices["powerbank"] = 200
	itemPrices["wallet"] = 50

	expected_first_user_info := handler.InfoResponse{
		Coins: INITIAL_MONEY - itemPrices["hoody"] - itemPrices["powerbank"] - TRANSACTION_MONEY,
		Inventory: []handler.Inventory{
			{
				Type:     "hoody",
				Quantity: 1,
			},
			{
				Type:     "powerbank",
				Quantity: 1,
			}},
		CoinHistory: handler.CoinHistory{
			Received: nil,
			Sent: []handler.SentTransaction{
				{
					ToUser: second_user,
					Amount: TRANSACTION_MONEY,
				}},
		},
	}

	expected_second_user_info := handler.InfoResponse{
		Coins: INITIAL_MONEY - itemPrices["umbrella"] - itemPrices["wallet"] + TRANSACTION_MONEY,
		Inventory: []handler.Inventory{
			{
				Type:     "umbrella",
				Quantity: 1,
			},
			{
				Type:     "wallet",
				Quantity: 1,
			}},
		CoinHistory: handler.CoinHistory{
			Received: []handler.ReceivedTransaction{
				{
					FromUser: first_user,
					Amount:   TRANSACTION_MONEY,
				}},
			Sent: nil,
		},
	}

	assert.Equal(t, expected_first_user_info, env.GetTestUserInfo())

	env.ShopClient.SetAuthToken(second_user_token.Token)
	code, second_user_info := env.ShopClient.GetInfo()
	assert.Equal(t, 200, code)
	assert.Equal(t, expected_second_user_info, second_user_info)
}
