package service

import (
	"time"
	"wishlist-bot/database"
)

type UserDto struct {
	ID        int64     `json:"id"`
	Name      string    `form:"name" json:"name"`
	Surname   string    `form:"surname" json:"surname"`
	Username  string    `form:"username" json:"patronymic"`
	Status    string    `form:"status" json:"status"`
	Birthdate time.Time `form:"birthdate" json:"birthdate"`
}

type UserService interface {
	Save(cRequest UserDto) bool
	FindById(ID int64) UserDto
	FindAllTotal() []UserDto
	FindAll(page, perPage int) ([]UserDto, *Pagination, error)
	UpdateBirthdate(birthdate *time.Time, ID int64) bool
	UpdateName(name string, ID int64) bool
	UpdateSurname(surname string, ID int64) bool
	UpdateUsername(username string, ID int64) bool
	UpdateStatus(status string, ID int64)
	Delete(ID int64)
	Restore(ID int64)
	ExistsById(ID int64) bool
	CheckIfDeleted(ID int64) bool
	CheckIfRegistered(ID int64) bool
}

type UserServiceImpl struct {
	Repo database.UserRepository
}

func (us *UserServiceImpl) Save(cRequest UserDto) bool {
	user := mapUserDtoToUser(&cRequest)
	return us.Repo.Save(user)
}

func (us *UserServiceImpl) FindById(id int64) UserDto {
	user := us.Repo.FindById(id)
	return *mapUserToDto(&user)
}

func (us *UserServiceImpl) FindAllTotal() []UserDto {
	users := us.Repo.FindAllTotal()
	var userDtos []UserDto
	for _, user := range users {
		dto := mapUserToDto(&user)
		userDtos = append(userDtos, *dto)
	}
	return userDtos
}

func (s *UserServiceImpl) FindAll(page, perPage int) ([]UserDto, *Pagination, error) {
	offset := (page - 1) * perPage
	users := s.Repo.FindAll(perPage, offset)
	var userDtos []UserDto
	for _, user := range users {
		dto := mapUserToDto(&user)
		userDtos = append(userDtos, *dto)
	}

	total := s.Repo.GetCount()

	pagination := NewPagination(total, perPage)
	pagination.CurrentPage = page

	return userDtos, pagination, nil
}

//func (us *UserServiceImpl) FindAll() []UserDto {
//	users := us.Repo.FindAll()
//	var result []UserDto
//	for _, u := range users {
//		dto := mapUserToDto(&u)
//		result = append(result, *dto)
//	}
//	return result
//}

func (us *UserServiceImpl) UpdateBirthdate(birthdate *time.Time, ID int64) bool {
	return us.Repo.UpdateBirthdate(*birthdate, ID)
}

func (us *UserServiceImpl) UpdateName(name string, ID int64) bool {
	return us.Repo.UpdateName(name, ID)
}

func (us *UserServiceImpl) UpdateSurname(surname string, ID int64) bool {
	return us.Repo.UpdateSurname(surname, ID)
}

func (us *UserServiceImpl) UpdateStatus(status string, ID int64) {
	us.Repo.UpdateStatus(status, ID)
}

func (us *UserServiceImpl) UpdateUsername(username string, ID int64) bool {
	return us.Repo.UpdateUsername(username, ID)
}

func (us *UserServiceImpl) Delete(id int64) {
	us.Repo.Delete(id)
}

func (us *UserServiceImpl) Restore(id int64) {
	us.Repo.Restore(id)
}

func (us *UserServiceImpl) ExistsById(id int64) bool {
	existsById := us.Repo.ExistsById(id)
	return existsById
}

func (us *UserServiceImpl) CheckIfDeleted(id int64) bool {
	return us.Repo.CheckIfDeleted(id)
}

func (us *UserServiceImpl) CheckIfRegistered(id int64) bool {
	return us.Repo.CheckIfRegistered(id)
}
