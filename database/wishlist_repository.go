package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type WRepository interface {
	Save(w *Wish) int64
	SaveAll(w []*Wish) error
	FindById(id int64) Wish
	FindAllByUserId(userId int64) []Wish
	Update(updateRequest Wish) int64
	Delete(id int64)
}

type WishlistRepository struct {
	DB *sql.DB
}

type Wish struct {
	ID        int64     `json:"id"`
	WishText  string    `json:"wish_text"`
	UserID    int64     `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (r *WishlistRepository) Save(w *Wish) int64 {
	query := `
	INSERT INTO wishes (wish_text, user_id)
	VALUES ($1, $2)
	RETURNING id;
`
	var id int64 = 0
	err := r.DB.QueryRow(query, &w.WishText, &w.UserID).Scan(&id)
	if err != nil {
		_ = r.DB.Close()
		fmt.Errorf("error at %s", err)
	}
	return id
}

func (r *WishlistRepository) SaveAll(wishes []*Wish) error {
	if len(wishes) == 0 {
		return errors.New("no wishes to save")
	}
	query := `INSERT INTO wishes(wish_text, user_id) VALUES `
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
		fmt.Errorf("error at %s", err)
		return err
	}
	return nil
}

func (r *WishlistRepository) FindById(ID int64) Wish {
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
		_ = r.DB.Close()
		fmt.Errorf("error at %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		errIn := rows.Scan(&w.ID, &w.WishText, &w.UserID)
		if errIn != nil {
			fmt.Errorf("Error at %s", errIn)
		}
	}
	return w
}

func (r *WishlistRepository) FindAllByUserId(ID int64) []Wish {
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
		_ = r.DB.Close()
		fmt.Errorf("error at %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		errIn := rows.Scan(&w.ID, &w.WishText, &w.UserID)
		if errIn != nil {
			fmt.Errorf("Error at %s", errIn)
		}
		wishes = append(wishes, w)
	}
	return wishes
}

func (r *WishlistRepository) Update(w Wish) int64 {
	query := `UPDATE wishes SET
		wish_text = $1
		WHERE user_id = $2
		AND deleted_at IS NULL
		RETURNING id`
	var id int64 = 0
	err := r.DB.QueryRow(query, &w.WishText, &w.UserID, &w.DeletedAt).Scan(&id)
	if err != nil {
		_ = r.DB.Close()
		fmt.Errorf("error at %s", err)
	}
	return id
}

func (r *WishlistRepository) Delete(ID int64) {
	query := `UPDATE wishes SET DELETED_at = NOW() WHERE id = $1`
	_, err := r.DB.Exec(query, &ID)
	if err != nil {
		_ = r.DB.Close()
		fmt.Errorf("error at %s", err)
	}
}
