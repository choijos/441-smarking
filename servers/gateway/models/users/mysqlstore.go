package users

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// SQLStore represents a mysql database containing all users
type SQLStore struct {
	DbStore *sql.DB
}

// GetByID takes in a user ID and returns the user from the database
// 	as a struct. Sends an error if user is not in the database or there
// 	are errors when scanning the database rows
func (ss *SQLStore) GetByID(id int64) (*User, error) {
	rows, err := ss.DbStore.Query("select id,email,passhash,username,first_name,last_name,photourl,phonenumber from users where id=?", id)

	if rows == nil || err != nil {
		return nil, ErrUserNotFound

	}

	defer rows.Close()

	getUser := User{}

	for rows.Next() {
		if err := rows.Scan(&getUser.ID, &getUser.Email, &getUser.PassHash,
			&getUser.UserName, &getUser.FirstName, &getUser.LastName, &getUser.PhotoURL, &getUser.PhoneNumber); err != nil {
			return nil, fmt.Errorf("error scanning row: %v\n", err)
		}

	}

	return &getUser, nil

}

// GetByEmail takes in a user email and returns the user from the database
// 	as a struct. Sends an error if user is not in the database or there
// 	are errors when scanning the database rows
func (ss *SQLStore) GetByEmail(email string) (*User, error) {
	rows, err := ss.DbStore.Query("select id,email,passhash,username,first_name,last_name,photourl,phonenumber from users where email=?", email)

	if rows == nil || err != nil {
		return nil, ErrUserNotFound

	}

	defer rows.Close()

	getUser := User{}

	for rows.Next() {
		if err := rows.Scan(&getUser.ID, &getUser.Email, &getUser.PassHash,
			&getUser.UserName, &getUser.FirstName, &getUser.LastName, &getUser.PhotoURL, &getUser.PhoneNumber); err != nil {
			return nil, fmt.Errorf("error scanning row: %v\n", err)
		}

	}

	return &getUser, nil

}

// GetByUserName takes in a user username and returns the user from the database
// 	as a struct. Sends an error if user is not in the database or there
// 	are errors when scanning the database rows
func (ss *SQLStore) GetByUserName(username string) (*User, error) {
	rows, err := ss.DbStore.Query("select id,email,passhash,username,first_name,last_name,photourl,phonenumber from users where username=?", username)

	if rows == nil || err != nil {
		return nil, ErrUserNotFound

	}

	defer rows.Close()

	getUser := User{}

	for rows.Next() {
		if err := rows.Scan(&getUser.ID, &getUser.Email, &getUser.PassHash,
			&getUser.UserName, &getUser.FirstName, &getUser.LastName, &getUser.PhotoURL, &getUser.PhoneNumber); err != nil {
			return nil, fmt.Errorf("error scanning row: %v\n", err)
		}

	}

	return &getUser, nil

}

// Insert takes in a user struct and adds this user to the database. Returns
// 	the user added with their new DBMS assigned ID. Returns error if user
// 	cannot be added to the database.
func (ss *SQLStore) Insert(user *User) (*User, error) {
	insq := "insert into users(email, first_name, last_name, username, passhash, photourl, phonenumber) values (?, ?, ?, ?, ?, ?, ?)" // regexp.QuoteMeta("insert into users(email, first_name, last_name, username, passhash, photourl) values (?, ?, ?, ?, ?, ?)")
	res, err := ss.DbStore.Exec(insq, user.Email, user.FirstName, user.LastName, user.UserName, user.PassHash, user.PhotoURL, user.PhoneNumber)
	if err != nil {
		return nil, ErrInvalidInsert

	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting new ID: %v\n", id)

	}

	user.ID = id

	return user, nil

}

// Update applies user updates to the user with the given user ID
// 	and returns the newly-updated user. Returns error if the user is
// 	not in the database
func (ss *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	insq := "update users set first_name = ?, last_name = ?, phonenumber = ? where id = ?"

	_, err := ss.DbStore.Exec(insq, updates.FirstName, updates.LastName, updates.PhoneNumber, id)
	if err != nil {
		return nil, ErrUserDNE

	}

	updatedUser, _ := ss.GetByID(id)

	return updatedUser, nil

}

// Delete removes the user with the given ID from the database.
// 	Returns error if user cannot be deleted.
func (ss *SQLStore) Delete(id int64) error {
	deleteString := "delete from users where id=?"
	_, err := ss.DbStore.Exec(deleteString, id)
	if err != nil {
		return fmt.Errorf("error deleting user with id %d: %v", id, err)

	}

	return nil

}

// InsertSignIn adds a user sign in attempt to the database. Returns an
// 	error if attempt cannot be inserted.
func (ss *SQLStore) InsertSignIn(id int64, currTime time.Time, ip string) error {
	insq := "insert into usersignin(id, whensignin, clientip) values (?, ?, ?)"
	_, err := ss.DbStore.Exec(insq, id, currTime, ip)
	if err != nil {
		return err

	}

	return nil

}