package wishlist

type Service struct {
	repo *Repository
}

func NewService(r *Repository) Service {
	return Service{repo: r}
}

func (w *Service) Save(wish Wish) error {
	return w.repo.Save(&wish)
}

func (w *Service) SaveAll(wishes []Wish) error {
	return w.repo.SaveAll(wishes)
}

func (w *Service) FindAllByUserId(userId int64) ([]Wish, error) {
	return w.repo.FindAllByUserId(userId)
}

func (w *Service) FindCountByUserID(userID int64) (int, error) {
	return w.repo.FindCountByUserID(userID)
}

func (w *Service) Update(wish Wish) error {
	return w.repo.Save(&wish)
}

func (w *Service) Delete(s string, userID int64) error {
	return w.repo.Delete(s, userID)
}

func (w *Service) DeleteAll(userID int64) error {
	return w.repo.DeleteAll(userID)
}

func (w *Service) DeleteByID(ID int64) error {
	return w.repo.DeleteByID(ID)
}
