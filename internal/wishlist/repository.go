package wishlist

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"wishlist-bot/internal/logger/sl"
)

type Repository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewRepository(db *sql.DB, log *slog.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func (r *Repository) Save(w *Wish) error {
	const op = "WishlistRepository.Save"

	query := `INSERT INTO wishes (wish_text, user_id) VALUES ($1, $2)`
	err := r.db.QueryRow(query, &w.WishText, &w.UserID)
	if err != nil {
		return fmt.Errorf("error at %v", err)
	}
	return nil
}

func (r *Repository) SaveAll(wishes []Wish) error {
	const op = "WishlistRepository.SaveAll"

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
	_, err := r.db.Exec(query, values...)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return fmt.Errorf("error at %s", err)
	}
	return nil
}

func (r *Repository) FindById(ID int64) (Wish, error) {
	const op = "WishlistRepository.FindById"

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
	rows, err := r.db.Query(query, ID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return Wish{}, fmt.Errorf("error at %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		errIn := rows.Scan(&w.ID, &w.WishText, &w.UserID)
		if errIn != nil {
			r.log.Error(op, sl.Err(err))
			return Wish{}, fmt.Errorf("error at %s", errIn)
		}
	}
	return w, nil
}

func (r *Repository) FindAllByUserId(ID int64) ([]Wish, error) {
	const op = "WishlistRepository.FindAllByUserId"

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
	rows, err := r.db.Query(query, ID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return nil, fmt.Errorf("error at %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&w.ID, &w.WishText, &w.UserID)
		if err != nil {
			r.log.Error(op, sl.Err(err))
			return nil, fmt.Errorf("error at %v", err)
		}
		wishes = append(wishes, w)
	}
	return wishes, nil
}

func (r *Repository) Delete(s string, userID int64) error {
	const op = "WishlistRepository.Delete"

	query := `DELETE FROM wishes WHERE wish_text LIKE $1 AND user_id = $2`
	_, err := r.db.Exec(query, &s, &userID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return fmt.Errorf("error at %s", err)
	}
	return nil
}

func (r *Repository) DeleteAll(userID int64) error {
	const op = "WishlistRepository.DeleteAll"

	query := `DELETE * FROM wishes WHERE user_id = $2`
	_, err := r.db.Exec(query, &userID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return errors.New("something went wrong with SQL")
	}
	return nil
}

func (r *Repository) DeleteByID(ID int64) error {
	const op = "WishlistRepository.DeleteByID"

	query := `DELETE * FROM wishes where id = $1`
	_, err := r.db.Exec(query, &ID)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return errors.New("something went wrong with SQL")
	}
	return nil
}
