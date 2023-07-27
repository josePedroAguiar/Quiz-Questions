package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int
	Username string
	Email    string
	Password string
	IsAdmin  bool
}

func AddUser(db *sql.DB, newUser User) string {
	// Create a new user
	
	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Store the hashed password in the database
	_, err = db.Exec("INSERT INTO tblUsers (username, email, password, is_admin) VALUES ($1, $2, $3, $4)",
		newUser.Username, newUser.Email, hashedPassword, newUser.IsAdmin)
	if err != nil {
		return "Email already exist!"
	}

	return "User added successfully!"
}

func getUserByEmailAndPassword(db *sql.DB, email, password string) *User {
	// Query the database to get the user by email
	row := db.QueryRow("SELECT user_id,username, email, password, is_admin FROM tblUsers WHERE email = $1", email)
	user := &User{}
	err := row.Scan(&user.ID,&user.Username, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			return nil
		}
		// Other query-related error
		return nil
	}

	// Verify the provided password against the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Password doesn't match
		return nil
	}

	// Password matches, return the user
	return user
}

func changeRoll(db *sql.DB, email string,is_admin string ) string {
	// Query the database to get the user by email
	
	updateStmt := ("Update tblUsers set is_admin=$1 WHERE email = $2")
	_, err := db.Exec(updateStmt, is_admin=="true",email )
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			return "User can't be updated"
		}
		// Other query-related error
		return "User can't be updated"
	}
	return "User roll updated"
}
