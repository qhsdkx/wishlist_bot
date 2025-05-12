package service

import db "wishlist-bot/database"

type WishDto struct {
	WishText string `form:"wish_text" json:"wish_text"`
	UserId   int64  `form:"user_id" json:"user_id"`
}

type WishService interface {
	Save(cRequest WishDto) int64
	SaveAll(wishList []WishDto) error
	FindAllByUserId(userId int64) []WishDto
	Update(uRequest WishDto) int64
	Delete(s string) error
}

type WishServiceImpl struct {
	Repository db.WRepository
}

func NewWishService(repository db.WRepository) WishService {
	return &WishServiceImpl{Repository: repository}
}

func (w *WishServiceImpl) Save(cRequest WishDto) int64 {
	wish := mapWishDtoToWish(cRequest)
	return w.Repository.Save(wish)
}

func (w *WishServiceImpl) SaveAll(wishList []WishDto) error {
	wishes := mapWishDtosToWishes(wishList)
	return w.Repository.SaveAll(wishes)
}

func (w *WishServiceImpl) FindAllByUserId(userId int64) []WishDto {
	wishes := w.Repository.FindAllByUserId(userId)
	var result []WishDto
	for _, wish := range wishes {
		dto := mapWishToWishDto(&wish)
		result = append(result, *dto)
	}
	return result
}

func (w *WishServiceImpl) Update(uRequest WishDto) int64 {
	wish := mapWishDtoToWish(uRequest)
	return w.Repository.Save(wish)
}

func (w *WishServiceImpl) Delete(s string) error {
	return w.Repository.Delete(s)
}
