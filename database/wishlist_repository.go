package database

//
//import (
//	"fmt"
//	"time"
//)
//
//type WRepository interface {
//	Save(w *Wish) int64
//	FindById(id int64) Wish
//	FindAllByUserId(userId int64) []Wish
//	Update(updateRequest Wish) int64
//	Delete(id int64)
//}
//
//type WishlistRepository struct{}
//
//type Wish struct {
//	ID        int64     `json:"id"`
//	WishText  string    `json:"wish_text"`
//	User      *User     `json:"user_id"`
//	DeletedAt time.Time `json:"deleted_at"`
//}
//
//func (*WishlistRepository) Save(w *Wish) int64 {
//	query := `
//	INSERT INTO wishes (wish_text, user_id)
//	VALUES ($1, $2)
//	RETURNING id;
//`
//	var id int64 = 0
//	err := db.QueryRow(query, &w.WishText, &w.User.ID).Scan(&id)
//	if err != nil {
//		_ = db.Close()
//		fmt.Errorf("error at %s", err)
//		fmt.Errorf("close connection to database")
//	}
//	return id
//}
//
//func (*WishlistRepository) FindById(ID int64) Wish {
//	query := `
//	SELECT
//    	w.id as id,
//    	w.wish_text as wishText,
//    	w.user_id as userId,
//		w.deleted_at
//	FROM wishes w
//	WHERE w.id = $1
//	AND deleted_at IS NULL
//`
//	w := Wish{}
//	rows, err := db.Query(query, ID)
//	if err != nil {
//		_ = db.Close()
//		fmt.Errorf("error at %s", err)
//		fmt.Errorf("close connection to database")
//	}
//	for rows.Next() {
//		errIn := rows.Scan(&w.ID, &w.WishText, &w.User.ID)
//		if errIn != nil {
//			fmt.Errorf("Error at %s", errIn)
//		}
//	}
//	return w
//}
//
//func (*WishlistRepository) FindAllByUserId(ID int64) []Wish {
//	var wishes []Wish
//	query := `
//	SELECT
//		w.id as id,
//		w.wish_text as wishText,
//		w.user_id as userId,
//		w.deleted_at as deletedAt
//	FROM wishes w
//	WHERE w.user_id = $1
//	AND deleted_at IS NULL
//`
//	w := Wish{}
//	rows, err := db.Query(query, ID)
//	if err != nil {
//		_ = db.Close()
//		fmt.Errorf("error at %s", err)
//		fmt.Errorf("close connection to database")
//	}
//	for rows.Next() {
//		errIn := rows.Scan(&w.ID, &w.WishText, &w.User.ID)
//		if errIn != nil {
//			fmt.Errorf("Error at %s", errIn)
//		}
//		wishes = append(wishes, w)
//	}
//	return wishes
//}
//
//func (*WishlistRepository) Update(w Wish) int64 {
//	query := `UPDATE wishes SET
//		wish_text = $1
//		WHERE user_id = $2
//		AND deleted_at IS NULL
//		RETURNING id`
//	var id int64 = 0
//	err := db.QueryRow(query, &w.WishText, &w.User.ID, &w.DeletedAt).Scan(&id)
//	if err != nil {
//		_ = db.Close()
//		fmt.Errorf("error at %s", err)
//		fmt.Errorf("close connection to database")
//	}
//	return id
//}
//
//func (*WishlistRepository) Delete(ID int64) {
//	query := `UPDATE wishes SET DELETED_at = NOW() WHERE id = $1`
//	_, err := db.Exec(query, &ID)
//	if err != nil {
//		_ = db.Close()
//		fmt.Errorf("error at %s", err)
//		fmt.Errorf("close connection to database")
//	}
//}
