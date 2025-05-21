package models

import (
	"database/sql"
	"errors"

	"strings"

	"golang.org/x/crypto/bcrypt"
	// "modernc.org/sqlite"
)

type User struct {
	Id int
	Name string
	Email string
	Password []byte
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) CreateUserAccount(name string, email string, 
	password string) (error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO Users (userName, userEmail, userPassword)
			VALUES (?, ?, ?);`

	result, err := um.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicatedEmail
		}
		return err
	}

	_,err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) Authenticate(email, password string) (int,error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT userId,  userPassword FROM Users
			WHERE userEmail = ?`
	
	err := um.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (um *UserModel) Exists(id int) (bool, error) {
	stmt := `SELECT userName, userEmail
			FROM Users
			WHERE userId = ?;`
	
	row := um.DB.QueryRow(stmt, id)
	var user User

	err := row.Scan(&user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return false, ErrNoRecord
		}
		return false, err
	}

	return true, nil
}

func (um *UserModel) ISUniqueEmail(email string) (bool) {
	var count int
	stmt := `SELECT COUNT(*)
			FROM Users 
			WHERE userEmail = ?;`
	
	row := um.DB.QueryRow(stmt, email)
	err := row.Scan(&count)
	if err != nil || count > 0 {
		return false
	}
	return true
}
