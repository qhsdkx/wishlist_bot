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
