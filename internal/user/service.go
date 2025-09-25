package user

import (
	"time"
	consta "wishlist-bot/internal/constant"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) FindByID(id int64) (User, error) {
	return s.repo.FindById(id)
}

func (us *Service) Save(user User) error {
	return us.repo.Save(&user)
}

func (us *Service) FindAllRegistered() ([]User, error) {
	return us.repo.FindAllTotal(consta.REGISTERED)
}

func (s *Service) FindAll(page, perPage int, mode string) ([]User, *Pagination, error) {
	offset := (page - 1) * perPage
	users, err := s.repo.FindAll(perPage, offset, mode)
	if err != nil {
		return nil, nil, err
	}

	total, err := s.repo.GetCount()
	if err != nil {
		return nil, nil, err
	}

	pagination := NewPagination(total, perPage)
	pagination.CurrentPage = page

	return users, pagination, nil
}

func (s *Service) FindAllUnregistered() ([]User, error) {
	return s.repo.FindAllTotal(consta.ADDED)
}

func (us *Service) UpdateBirthdate(birthdate *time.Time, ID int64) error {
	return us.repo.UpdateBirthdate(*birthdate, ID)
}

func (us *Service) UpdateName(name string, ID int64) error {
	return us.repo.UpdateName(name, ID)
}

func (us *Service) UpdateSurname(surname string, ID int64) error {
	return us.repo.UpdateSurname(surname, ID)
}

func (us *Service) UpdateStatus(status string, ID int64) error {
	return us.repo.UpdateStatus(status, ID)
}

func (us *Service) UpdateUsername(username string, ID int64) error {
	return us.repo.UpdateUsername(username, ID)
}

func (us *Service) Delete(id int64) error {
	return us.repo.Delete(id)
}

func (us *Service) ExistsById(id int64) error {
	return us.repo.ExistsById(id)
}

func (us *Service) CheckIfRegistered(id int64) error {
	return us.repo.CheckIfRegistered(id)
}
