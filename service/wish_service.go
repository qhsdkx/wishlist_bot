package service

import db "wishlist-bot/database"

type WishDto struct {
	WishText string `form:"wish_text" json:"wish_text"`
	UserId   int64  `form:"user_id" json:"user_id"`
}

type WishService interface {
	Save(cRequest WishDto) error
	SaveAll(wishList []WishDto) error
	FindAllByUserId(userId int64) ([]WishDto, error)
	Delete(s string, userID int64) error
	DeleteAll(userID int64) error
}

type WishServiceImpl struct {
	Repository db.WRepository
}

func NewWishService(repository db.WRepository) WishService {
	return &WishServiceImpl{Repository: repository}
}

func (w *WishServiceImpl) Save(cRequest WishDto) error {
	wish := mapWishDtoToWish(cRequest)
	return w.Repository.Save(wish)
}

func (w *WishServiceImpl) SaveAll(wishList []WishDto) error {
	wishes := mapWishDtosToWishes(wishList)
	return w.Repository.SaveAll(wishes)
}

func (w *WishServiceImpl) FindAllByUserId(userId int64) ([]WishDto, error) {
	wishes, err := w.Repository.FindAllByUserId(userId)
	if err != nil {
		return nil, err
	}
	result := make([]WishDto, len(wishes))
	for i, wish := range wishes {
		dto := mapWishToWishDto(&wish)
		result[i] = *dto
	}
	return result, nil
}

func (w *WishServiceImpl) Update(uRequest WishDto) error {
	wish := mapWishDtoToWish(uRequest)
	return w.Repository.Save(wish)
}

func (w *WishServiceImpl) Delete(s string, userID int64) error {
	return w.Repository.Delete(s, userID)
}

func (w *WishServiceImpl) DeleteAll(userID int64) error {
	return w.Repository.DeleteAll(userID)
}
