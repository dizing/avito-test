package domain

// import "context"

// type ShopService struct {
// 	itemRepo ItemRepository
// 	userRepo UserRepository
// }

// func NewService(itemRepo ItemRepository, userRepo UserRepository) *ShopService {
// 	return &ShopService{
// 		userRepo: userRepo,
// 		itemRepo: itemRepo,
// 	}
// }

// func (s *ShopService) BuyItem(context context.Context, user_id UserUUID, item_name string) error {
// 	user, _ := s.userRepo.GetById(context, user_id)

// 	item, _ := s.itemRepo.GetItemByName(context, item_name)

// 	return nil
// }
