package service

import db "wishlist-bot/database"

func mapUserDtoToUser(dto *UserDto) *db.User {
	return &db.User{ID: dto.ID, Name: dto.Name, Surname: dto.Surname, Username: "@" + dto.Username, Status: dto.Status, Birthdate: dto.Birthdate}
}

func mapUserToDto(u *db.User) *UserDto {
	return &UserDto{u.ID, u.Name, u.Surname, u.Username, u.Status, u.Birthdate}
}

func mapWishToWishDto(wish *db.Wish) *WishDto {
	return &WishDto{WishText: wish.WishText, UserId: wish.UserID}
}

func mapWishDtoToWish(wishDto WishDto) *db.Wish {
	return &db.Wish{WishText: wishDto.WishText, UserID: wishDto.UserId}
}

func mapWishDtosToWishes(wishes []WishDto) []*db.Wish {
	result := make([]*db.Wish, len(wishes))
	for i, wish := range wishes {
		wishDto := mapWishDtoToWish(wish)
		result[i] = wishDto
	}
	return result
}
