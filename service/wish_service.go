package service

//
//import db "wishlist-bot/database"
//
//type WishDto struct {
//	WishText string `form:"wish_text" json:"wish_text"`
//	UserId   int64  `form:"user_id" json:"user_id"`
//}
//
//type WishService interface {
//	Save(cRequest WishDto) int64
//	FindAllByUserId(userId int64) []WishDto
//	Update(uRequest WishDto) int64
//	Delete(id int64)
//}
//
//var wRepository = db.WishlistRepository{}
//
//type WishServiceImpl struct{}
//
//func (w *WishServiceImpl) Save(cRequest WishDto) int64 {
//	wish := mapWishDtoToWish(&cRequest)
//	return wRepository.Save(wish)
//}
//
//func (w *WishServiceImpl) FindAllByUserId(userId int64) []WishDto {
//	wishes := wRepository.FindAllByUserId(userId)
//	var result []WishDto
//	for _, wish := range wishes {
//		dto := mapWishToWishDto(&wish)
//		result = append(result, *dto)
//	}
//	return result
//}
//
//func (w *WishServiceImpl) Update(uRequest WishDto) int64 {
//	wish := mapWishDtoToWish(&uRequest)
//	return wRepository.Save(wish)
//}
//
//func (w *WishServiceImpl) Delete(id int64) {
//	wRepository.Delete(id)
//}
