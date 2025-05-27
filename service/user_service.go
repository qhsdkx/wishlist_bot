package service

import (
	"time"
	consta "wishlist-bot/constant"
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
	Save(cRequest UserDto) error
	FindById(ID int64) (UserDto, error)
	FindAllRegistered() ([]UserDto, error)
	FindAll(page, perPage int) ([]UserDto, *Pagination, error)
	FindAllUnregistered() ([]UserDto, error)
	UpdateBirthdate(birthdate *time.Time, ID int64) error
	UpdateName(name string, ID int64) error
	UpdateSurname(surname string, ID int64) error
	UpdateUsername(username string, ID int64) error
	UpdateStatus(status string, ID int64) error
	Delete(ID int64) error
	ExistsById(ID int64) error
	CheckIfRegistered(ID int64) error
}

type UserServiceImpl struct {
	Repo database.UserRepository
}

func NewUserService(repo database.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

func (us *UserServiceImpl) Save(cRequest UserDto) error {
	user := mapUserDtoToUser(&cRequest)
	return us.Repo.Save(user)
}

func (us *UserServiceImpl) FindById(id int64) (UserDto, error) {
	user, err := us.Repo.FindById(id)
	if err != nil {
		return UserDto{}, err
	}
	return *mapUserToDto(&user), nil
}

func (us *UserServiceImpl) FindAllRegistered() ([]UserDto, error) {
	users, err := us.Repo.FindAllTotal(consta.REGISTERED)
	if err != nil {
		return nil, err
	}
	userDtos := make([]UserDto, len(users))
	for i, user := range users {
		dto := mapUserToDto(&user)
		userDtos[i] = *dto
	}
	return userDtos, nil
}

func (s *UserServiceImpl) FindAll(page, perPage int) ([]UserDto, *Pagination, error) {
	offset := (page - 1) * perPage
	users, err := s.Repo.FindAll(perPage, offset)
	if err != nil {
		return nil, nil, err
	}
	userDtos := make([]UserDto, len(users))
	for i, user := range users {
		dto := mapUserToDto(&user)
		userDtos[i] = *dto
	}

	total, err := s.Repo.GetCount()
	if err != nil {
		return nil, nil, err
	}

	pagination := NewPagination(total, perPage)
	pagination.CurrentPage = page

	return userDtos, pagination, nil
}

func (s *UserServiceImpl) FindAllUnregistered() ([]UserDto, error) {
	users, err := s.Repo.FindAllTotal(consta.ADDED)
	if err != nil {
		return nil, err
	}
	userDtos := make([]UserDto, len(users))
	for i, user := range users {
		dto := mapUserToDto(&user)
		userDtos[i] = *dto
	}
	return userDtos, nil
}

func (us *UserServiceImpl) UpdateBirthdate(birthdate *time.Time, ID int64) error {
	return us.Repo.UpdateBirthdate(*birthdate, ID)
}

func (us *UserServiceImpl) UpdateName(name string, ID int64) error {
	return us.Repo.UpdateName(name, ID)
}

func (us *UserServiceImpl) UpdateSurname(surname string, ID int64) error {
	return us.Repo.UpdateSurname(surname, ID)
}

func (us *UserServiceImpl) UpdateStatus(status string, ID int64) error {
	return us.Repo.UpdateStatus(status, ID)
}

func (us *UserServiceImpl) UpdateUsername(username string, ID int64) error {
	return us.Repo.UpdateUsername(username, ID)
}

func (us *UserServiceImpl) Delete(id int64) error {
	return us.Repo.Delete(id)
}

func (us *UserServiceImpl) ExistsById(id int64) error {
	return us.Repo.ExistsById(id)
}

func (us *UserServiceImpl) CheckIfRegistered(id int64) error {
	return us.Repo.CheckIfRegistered(id)
}
