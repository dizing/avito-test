package domain

import (
	"github.com/google/uuid"
)

type ItemName string

type Item struct {
	Name  ItemName
	Price uint
}

type TransactionUUID uuid.UUID
type UserName string

type Transaction struct {
	Id     TransactionUUID
	From   UserName
	To     UserName
	Amount uint
}

type TransactionHistory []*Transaction

func NewTransaction(From UserName, To UserName, amount uint) *Transaction {
	return &Transaction{
		Id:     TransactionUUID{},
		From:   From,
		To:     To,
		Amount: amount,
	}
}

type Possession struct {
	Username UserName
	Item     ItemName
	Amount   uint
}

type UserInventory map[ItemName]uint

type User struct {
	Username     UserName
	PasswordHash []byte
	Balance      int
}

func NewUser(username UserName, password []byte) *User {
	return &User{
		Username:     username,
		PasswordHash: password,
		Balance:      5000,
	}
}

type UserInfo struct {
	User         *User
	Inventory    UserInventory
	Transactions TransactionHistory
}

func NewUserInventory(user_posessions []*Possession) *UserInventory {
	var inventory = UserInventory{}

	for _, possession := range user_posessions {
		inventory[possession.Item] = possession.Amount
	}

	return &inventory
}
