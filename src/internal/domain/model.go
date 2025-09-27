package domain

import (
	"github.com/google/uuid"
)

type ItemName string

type Item struct {
	Name  ItemName
	Price int
}

type TransactionUUID uuid.UUID
type UserName string

type Transaction struct {
	Id     TransactionUUID
	From   UserName
	To     UserName
	Amount int
}

type TransactionHistory []*Transaction

func NewTransaction(From UserName, To UserName, amount int) *Transaction {
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
	Amount   int
}

func NewEmptyPossession(username UserName, item ItemName) *Possession {
	return &Possession{
		Username: username,
		Item:     item,
		Amount:   0,
	}
}

type UserInventory map[ItemName]int

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
