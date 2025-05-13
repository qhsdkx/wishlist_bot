package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type WRepository interface {
	Save(w *Wish) error
	SaveAll(w []*Wish) error
	FindById(id int64) (Wish, error)
	FindAllByUserId(userId int64) ([]Wish, error)
	Update(updateRequest Wish) error
	Delete(s string) error
}

type WishlistRepository struct {
	DB *sql.DB
}

func NewWishlistRepository(db *sql.DB) *WishlistRepository {
	return &WishlistRepository{DB: db}
}

type Wish struct {
	ID        int64     `json:"id"`
	WishText  string    `json:"wish_text"`
	UserID    int64     `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (r *WishlistRepository) Save(w *Wish) error {
	query := `INSERT INTO wishes (wish_text, user_id) VALUES ($1, $2)`
	err := r.DB.QueryRow(query, &w.WishText, &w.UserID)
	if err != nil {
		return fmt.Errorf("error at %s", err)
	}
	return nil
}

func (r *WishlistRepository) SaveAll(wishes []*Wish) error {
	if len(wishes) == 0 {
		return errors.New("no wishes to save")
	}
	query := `INSERT INTO wishes(wish_text, user_id) VALUES`
	values := make([]interface{}, 0, len(wishes)*2)
	for i, w := range wishes {
		if i > 0 {
			query += ", "
		}
		query += "($" + strconv.Itoa(i*2+1) + ", $" + strconv.Itoa(i*2+2) + ")"
		values = append(values, w.WishText, w.UserID)
	}
	_, err := r.DB.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("error at %s", err)
	}
	return nil
}

func (r *WishlistRepository) FindById(ID int64) (Wish, error) {
	query := `
	SELECT
   	w.id as id,
   	w.wish_text as wishText,
   	w.user_id as userId,
		w.deleted_at
	FROM wishes w
	WHERE w.id = $1
	AND deleted_at IS NULL
`
	w := Wish{}
	rows, err := r.DB.Query(query, ID)
	if err != nil {
		return Wish{}, fmt.Errorf("error at %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		errIn := rows.Scan(&w.ID, &w.WishText, &w.UserID)
		if errIn != nil {
			return Wish{}, fmt.Errorf("error at %s", errIn)
		}
	}
	return w, nil
}

func (r *WishlistRepository) FindAllByUserId(ID int64) ([]Wish, error) {
	var wishes []Wish
	query := `
	SELECT
		w.id as id,
		w.wish_text as wishText,
		w.user_id as userId
	FROM wishes w
	WHERE w.user_id = $1
	AND deleted_at IS NULL
`
	w := Wish{}
	rows, err := r.DB.Query(query, ID)
	if err != nil {
		return nil, fmt.Errorf("error at %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		errIn := rows.Scan(&w.ID, &w.WishText, &w.UserID)
		if errIn != nil {
			return nil, fmt.Errorf("Error at %s", errIn)
		}
		wishes = append(wishes, w)
	}
	return wishes, nil
}

func (r *WishlistRepository) Update(w Wish) error {
	query := `UPDATE wishes SET
		wish_text = $1
		WHERE user_id = $2
		AND deleted_at IS NULL`
	err := r.DB.QueryRow(query, &w.WishText, &w.UserID, &w.DeletedAt)
	if err != nil {
		return fmt.Errorf("error at %s", err)
	}
	return nil
}

func (r *WishlistRepository) Delete(s string) error {
	var deletedAt time.Time
	query := `UPDATE wishes SET DELETED_AT = NOW() WHERE wish_text LIKE $1 RETURNING deleted_at`
	err := r.DB.QueryRow(query, &s).Scan(&deletedAt)
	if err != nil || deletedAt.IsZero() {
		return fmt.Errorf("error at %s", err)
	}
	return nil
}
