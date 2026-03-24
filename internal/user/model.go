package user

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	Birthdate time.Time `json:"birthdate"`
}

func (u *User) GetID() int64 {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetSurname() string {
	return u.Surname
}
