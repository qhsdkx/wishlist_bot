package database

import (
	"database/sql"
	"fmt"
	"time"
	constants "wishlist-bot/constant"
)

type UserRepository interface {
	Save(user *User) error
	FindById(id int64) (User, error)
	FindAllTotal() ([]User, error)
	FindAll(page, perPage int) ([]User, error)
	UpdateBirthdate(birthdate time.Time, ID int64) error
	UpdateName(name string, ID int64) error
	UpdateSurname(surname string, ID int64) error
	UpdateUsername(username string, ID int64) error
	UpdateStatus(status string, ID int64) error
	Delete(id int64) error
	Restore(ID int64) error
	ExistsById(id int64) error
	CheckIfDeleted(ID int64) error
	CheckIfRegistered(ID int64) error
	GetCount() (int, error)
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	Birthdate time.Time `json:"birthdate"`
}

func (ur *UserRepositoryImpl) Save(u *User) error {
	query := `INSERT INTO users(id, name, surname, username, birthdate, status)
			  VALUES($1, $2, $3, $4, $5, $6)`
	_, err := ur.DB.Exec(query, &u.ID, &u.Name, &u.Surname, &u.Username, &u.Birthdate, constants.ADDED)
	return err
}

func (ur *UserRepositoryImpl) FindById(ID int64) (User, error) {
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
	rows, err := ur.DB.Query(query, ID)
	if err != nil {
		return u, fmt.Errorf("error at %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		errIn := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Birthdate, &u.Username, &u.Status)
		if errIn != nil {
			fmt.Errorf("error at %s", errIn)
		}
	}
	return u, nil
}

func (ur *UserRepositoryImpl) FindAllTotal() ([]User, error) {
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
		AND u.status = $1
	`
	var users []User
	rows, err := ur.DB.Query(query, constants.REGISTERED)
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

func (ur *UserRepositoryImpl) FindAll(page, perPage int) ([]User, error) {
	var users []User

	query := `SELECT 
    	u.id as id,
	   	u.name as name,
	   	u.surname as surname,
	   	u.birthdate as birthdate,
	   	u.username as username,
	   	u.status as status
		FROM users u
        WHERE status = $1 
        ORDER BY name 
        LIMIT $2 OFFSET $3`
	rows, err := ur.DB.Query(query, constants.REGISTERED, page, perPage)
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

func (ur *UserRepositoryImpl) UpdateBirthdate(birthdate time.Time, ID int64) error {
	query := `UPDATE users
	SET
    birthdate = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.DB.Exec(query, &birthdate, &ID)
	if err != nil {
		return fmt.Errorf("error at %s", err)
	}
	return nil
}

func (ur *UserRepositoryImpl) UpdateName(name string, ID int64) error {
	query := `UPDATE users
	SET
    name = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.DB.Exec(query, &name, &ID)
	if err != nil {
		return fmt.Errorf("error at Saving name")
	}
	return nil
}

func (ur *UserRepositoryImpl) UpdateSurname(surname string, ID int64) error {
	query := `UPDATE users
SET
    surname = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.DB.Exec(query, &surname, &ID)
	if err != nil {
		return fmt.Errorf("error at saving surname")
	}
	return nil
}

func (ur *UserRepositoryImpl) UpdateUsername(username string, ID int64) error {
	query := `UPDATE users
	SET
    username = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.DB.Exec(query, &username, &ID)
	if err != nil {
		return fmt.Errorf("error at saving username")
	}
	return nil
}

func (ur *UserRepositoryImpl) UpdateStatus(status string, ID int64) error {
	query := `UPDATE users
	SET
    status = $1,
    updated_at = now()
    WHERE deleted_at IS NULL
    AND id = $2`
	_, err := ur.DB.Exec(query, &status, &ID)
	if err != nil {
		return fmt.Errorf("error at update status")
	}
	return nil
}

func (ur *UserRepositoryImpl) Delete(ID int64) error {
	query := `UPDATE users
	SET deleted_at = now(),
	    updated_at = now()
	WHERE id = $1`
	_, err := ur.DB.Exec(query, ID)
	if err != nil {
		return fmt.Errorf("error at delete user")
	}
	return nil
}

func (ur *UserRepositoryImpl) Restore(ID int64) error {
	query := `UPDATE users
	SET deleted_at = NULL,
	    updated_at = now()
	WHERE id = $1`
	_, err := ur.DB.Exec(query, ID)
	if err != nil {
		return fmt.Errorf("error at restore user")
	}
	return nil
}

func (ur *UserRepositoryImpl) ExistsById(ID int64) error {
	var result bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	exists, err := ur.DB.Query(query, &ID)
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

func (ur *UserRepositoryImpl) CheckIfDeleted(ID int64) error {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NOT NULL)`
	exists, err := ur.DB.Query(query, &ID)
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
	if result {
		return fmt.Errorf("user is deleted")
	}
	return nil
}

func (ur *UserRepositoryImpl) CheckIfRegistered(ID int64) error {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND status = $2)`
	exists, err := ur.DB.Query(query, &ID, constants.REGISTERED)
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

func (ur *UserRepositoryImpl) GetCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE status = $1`
	rows, err := ur.DB.Query(query, constants.REGISTERED)
	if err != nil {
		fmt.Errorf("error at %s", err)
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
