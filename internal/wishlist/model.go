package wishlist

import (
	"time"
)

type Wish struct {
	ID        int64     `json:"id"`
	WishText  string    `json:"wish_text"`
	UserID    int64     `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
