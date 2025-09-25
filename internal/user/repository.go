package user

import (
	"database/sql"
	"fmt"
	"time"
	constants "wishlist-bot/internal/constant"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (ur *Repository) Save(u *User) error {
	query := `INSERT INTO users(id, name, surname, username, birthdate, status)
			  VALUES($1, $2, $3, $4, $5, $6)`
	_, err := ur.db.Exec(query, &u.ID, &u.Name, &u.Surname, &u.Username, &u.Birthdate, constants.ADDED)
	return err
}

func (ur *Repository) FindById(ID int64) (User, error) {
	query := `
	SELECT 
    	u.id as id,
    	u.name as name,
    	u.surname as surname,
    	u.birthdate as birthdate,
    	u.username as username,
    	u.status as status
	FROM users u
	WHERE u.id = $1
	AND u.deleted_at IS NULL
`
	u := User{}
	rows, err := ur.db.Query(query, ID)
	if err != nil {
		return u, fmt.Errorf("error at %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		errIn := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Birthdate, &u.Username, &u.Status)
		if errIn != nil {
			return User{}, fmt.Errorf("error at %s", errIn)
		}
	}
	return u, nil
}

func (ur *Repository) FindAllTotal(status string) ([]User, error) {
	query := `
		SELECT
	   	u.id as id,
	   	u.name as name,
	   	u.surname as surname,
	   	u.birthdate as birthdate,
	   	u.username as username,
	   	u.status as status
		FROM users u
		WHERE u.deleted_at IS NULL
		AND CASE WHEN $1 != 'N' THEN u.status = $1 ELSE TRUE END;
	`
	var users []User
	rows, err := ur.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("error at %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		errIn := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Birthdate, &u.Username, &u.Status)
		if errIn != nil {
			return nil, fmt.Errorf("error at %s", errIn)
		}
		users = append(users, u)
	}
	return users, nil
}

func (ur *Repository) FindAll(page, perPage int, mode string) ([]User, error) {
	var users []User

	query := `SELECT 
    	u.id as id,
	   	u.name as name,
	   	u.surname as surname,
	   	u.birthdate as birthdate,
	   	u.username as username,
	   	u.status as status
		FROM users u
        WHERE CASE WHEN ($1 != 'MIS') THEN u.status = $1 ELSE TRUE END
        ORDER BY name 
        LIMIT $2 OFFSET $3`
	var rows *sql.Rows
	var err error
	if mode == constants.SHOW_USERS {
		rows, err = ur.db.Query(query, constants.REGISTERED, page, perPage)
		if err != nil {
			return nil, fmt.Errorf("error at %s", err)
		}
	}
	if mode == constants.SEND_MESSAGE_ADMIN {
		rows, err = ur.db.Query(query, "MIS", page, perPage)
		if err != nil {
			return nil, fmt.Errorf("error at %s", err)
		}
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		errIn := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Birthdate, &u.Username, &u.Status)
		if errIn != nil {
			return nil, fmt.Errorf("error at %s", errIn)
		}
		users = append(users, u)
	}
	return users, nil
}

func (ur *Repository) UpdateBirthdate(birthdate time.Time, ID int64) error {
	query := `UPDATE users
	SET
    birthdate = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.db.Exec(query, &birthdate, &ID)
	if err != nil {
		return fmt.Errorf("error at %s", err)
	}
	return nil
}

func (ur *Repository) UpdateName(name string, ID int64) error {
	query := `UPDATE users
	SET
    name = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.db.Exec(query, &name, &ID)
	if err != nil {
		return fmt.Errorf("error at Saving name")
	}
	return nil
}

func (ur *Repository) UpdateSurname(surname string, ID int64) error {
	query := `UPDATE users
	SET
    surname = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.db.Exec(query, &surname, &ID)
	if err != nil {
		return fmt.Errorf("error at saving surname")
	}
	return nil
}

func (ur *Repository) UpdateUsername(username string, ID int64) error {
	query := `UPDATE users
	SET
    username = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.db.Exec(query, &username, &ID)
	if err != nil {
		return fmt.Errorf("error at saving username")
	}
	return nil
}

func (ur *Repository) UpdateStatus(status string, ID int64) error {
	query := `UPDATE users
	SET
    status = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.db.Exec(query, &status, &ID)
	if err != nil {
		return fmt.Errorf("error at update status")
	}
	return nil
}

func (ur *Repository) Delete(ID int64) error {
	wQuery := `DELETE FROM wishes WHERE user_id = $1`
	_, err := ur.db.Exec(wQuery, ID)
	if err != nil {
		return fmt.Errorf("error at delete wishes of user with id %d", ID)
	}
	query := `DELETE FROM users WHERE id = $1`
	_, err = ur.db.Exec(query, ID)
	if err != nil {
		return fmt.Errorf("error at delete user")
	}
	return nil
}

func (ur *Repository) ExistsById(ID int64) error {
	var result bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	exists, err := ur.db.Query(query, &ID)
	if err != nil {
		return fmt.Errorf("error at %s", err)
	}

	defer exists.Close()

	for exists.Next() {
		errRead := exists.Scan(&result)
		if errRead != nil {
			return fmt.Errorf("error at %s", err)
		}
	}
	if !result {
		return fmt.Errorf("user doesn't exist in database")
	}
	return nil
}

func (ur *Repository) CheckIfRegistered(ID int64) error {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND status = $2)`
	exists, err := ur.db.Query(query, &ID, constants.REGISTERED)
	if err != nil {
		return fmt.Errorf("error at %s", err)
	}

	defer exists.Close()

	if exists.Next() {
		errRead := exists.Scan(&result)
		if errRead != nil {
			return fmt.Errorf("error at %s", err)
		}
	}
	if !result {
		return fmt.Errorf("user is unregistered")
	}
	return nil
}

func (ur *Repository) GetCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE status = $1`
	rows, err := ur.db.Query(query, constants.REGISTERED)
	if err != nil {
		return 0, fmt.Errorf("error at %s", err)
	}

	defer rows.Close()

	if rows.Next() {
		errRead := rows.Scan(&count)
		if errRead != nil {
			return 0, fmt.Errorf("error at %s", err)
		}
	}
	return count, nil
}
