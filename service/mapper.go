package service

import db "wishlist-bot/database"

func mapUserDtoToUser(dto *UserDto) *db.User {
	return &db.User{ID: dto.ID, Name: dto.Name, Surname: dto.Surname, Username: "@" + dto.Username, Status: dto.Status, Birthdate: dto.Birthdate}
}

func mapUserToDto(u *db.User) *UserDto {
	return &UserDto{u.ID, u.Name, u.Surname, u.Username, u.Status, u.Birthdate}
}

//func mapWishToWishDto(wish *db.Wish) *WishDto {
//	return &WishDto{WishText: wish.WishText, UserId: wish.User.ID}
//}
//
//func mapWishDtoToWish(wishDto *WishDto) *db.Wish {
//	return &db.Wish{WishText: wishDto.WishText,
//		User: &db.User{ID: wishDto.UserId}}
//}
