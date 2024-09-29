package dbtools

import (
	"database/sql"
	"fmt"
)

type User struct {
	User_id  int64  `field:"user_id"`
	Username string `field:"username"`
	Email    string `field:"email"`
	Password string `field:"password"`
}

func (u *User) AddToDB(db *sql.DB) error {
	// Make statement
	prompt := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	statement, err := db.Prepare(prompt)
	if err != nil {
		fmt.Println("Error at statement")
		return err
	}

	// Execute statement
	result, err := statement.Exec(u.Username, u.Email, u.Password)
	if err != nil {
		fmt.Println("Error at execution")
		return err
	}

	// Give user there id
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error at id")
		return err
	}
	fmt.Println(id)
	u.User_id = id

	return nil
}

func GetUserFromEmail(email string, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE email=?"
	row := db.QueryRow(prompt, email)

	user := EmptyUser()
	err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Password)

	return user, err
}

func GetUserFromUsername(username string, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE username=?"
	row := db.QueryRow(prompt, username)

	user := EmptyUser()
	err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Password)

	return user, err
}

func NewUser(username, email, password string) *User {
	return &User{
		User_id:  -1,
		Username: username,
		Email:    email,
		Password: password,
	}
}

func EmptyUser() *User {
	return &User{User_id: -1}
}
