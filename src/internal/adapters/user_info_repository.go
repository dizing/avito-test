package adapters

import (
	"avito-test/internal/domain"
	"context"
)

type userInfoRepository struct {
	userRepo        domain.UserRepository
	possessionRepo  domain.PossessionRepository
	transactionRepo domain.TransactionRepository
}

func NewInfoRepository(userRepo domain.UserRepository, possessionRepo domain.PossessionRepository, transactionRepo domain.TransactionRepository) domain.UserInfoRepository {
	return &userInfoRepository{
		userRepo:        userRepo,
		possessionRepo:  possessionRepo,
		transactionRepo: transactionRepo,
	}
}

func (r *userInfoRepository) GetByUsername(ctx context.Context, username domain.UserName) (*domain.UserInfo, error) {
	user, err := r.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	userInventory, err := r.possessionRepo.GetUserInventory(ctx, username)
	if err != nil {
		return nil, err
	}

	transactionHistory, err := r.transactionRepo.GetAllByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		User:         user,
		Inventory:    userInventory,
		Transactions: transactionHistory,
	}, nil
}
