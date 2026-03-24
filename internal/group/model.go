package group

import "time"

const (
	GroupStatusUpcoming = "upcoming"
	GroupStatusPassed   = "passed"
)

type Group struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	BirthdayUserID int64     `json:"birthday_user_id"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type GroupMember struct {
	ID       int64     `json:"id"`
	GroupID  int64     `json:"group_id"`
	UserID   int64     `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}
