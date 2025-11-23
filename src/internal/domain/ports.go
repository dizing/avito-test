package domain

import (
	"context"
	"fmt"
)

var ErrEntityDoesNotExist = fmt.Errorf("entity does not exist")

type UserRepository interface {
	GetByUsername(context.Context, UserName) (*User, error)
	Save(context.Context, *User) error
}

type UserInfoRepository interface {
	GetByUsername(context.Context, UserName) (*UserInfo, error)
}

type TransactionRepository interface {
	GetById(context.Context, TransactionUUID) (*Transaction, error)
	GetTransactionHistoryByUsername(context.Context, UserName) (TransactionHistory, error)
	Save(context.Context, *Transaction) error
}

type ItemRepository interface {
	GetByName(context.Context, ItemName) (*Item, error)
}

type PossessionRepository interface {
	GetPossessionByUsernameAndItemName(context.Context, UserName, ItemName) (*Possession, error)
	GetUserInventory(context.Context, UserName) (UserInventory, error)
	Save(context.Context, *Possession) error
}
