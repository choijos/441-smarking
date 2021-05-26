package users

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SQLStore struct {
	DbStore *sql.DB
}

func (ss *SQLStore) GetByID(id int64) (*User, error) {
	//select rows from the table
	rows, err := ss.DbStore.Query("select id,email,passhash,username,first_name,last_name,photourl from users where id=?", id)

	if rows == nil || err != nil {
		return nil, ErrUserNotFound

	}

	defer rows.Close()

	getUser := User{}

	//while there are more rows
	for rows.Next() {
		if err := rows.Scan(&getUser.ID, &getUser.Email, &getUser.PassHash,
			&getUser.UserName, &getUser.FirstName, &getUser.LastName, &getUser.PhotoURL); err != nil {
			return nil, fmt.Errorf("error scanning row: %v\n", err)
		}

	}

	return &getUser, nil

}

func (ss *SQLStore) GetByEmail(email string) (*User, error) {
	//select rows from the table
	rows, err := ss.DbStore.Query("select id,email,passhash,username,first_name,last_name,photourl from users where email=?", email)

	if rows == nil || err != nil {
		return nil, ErrUserNotFound

	}

	defer rows.Close()

	getUser := User{}

	//while there are more rows
	for rows.Next() {
		if err := rows.Scan(&getUser.ID, &getUser.Email, &getUser.PassHash,
			&getUser.UserName, &getUser.FirstName, &getUser.LastName, &getUser.PhotoURL); err != nil {
			return nil, fmt.Errorf("error scanning row: %v\n", err)
		}

	}

	return &getUser, nil

}

func (ss *SQLStore) GetByUserName(username string) (*User, error) {
	//select rows from the table
	rows, err := ss.DbStore.Query("select id,email,passhash,username,first_name,last_name,photourl from users where username=?", username)

	if rows == nil || err != nil {
		return nil, ErrUserNotFound

	}

	defer rows.Close()

	getUser := User{}

	//while there are more rows
	for rows.Next() {
		if err := rows.Scan(&getUser.ID, &getUser.Email, &getUser.PassHash,
			&getUser.UserName, &getUser.FirstName, &getUser.LastName, &getUser.PhotoURL); err != nil {
			return nil, fmt.Errorf("error scanning row: %v\n", err)
		}

	}

	return &getUser, nil

}

func (ss *SQLStore) Insert(user *User) (*User, error) {
	//insert a new row into the table
	//use ? markers for the values to defeat SQL
	//injection attacks
	insq := "insert into users(email, first_name, last_name, username, passhash, photourl) values (?, ?, ?, ?, ?, ?)" // regexp.QuoteMeta("insert into users(email, first_name, last_name, username, passhash, photourl) values (?, ?, ?, ?, ?, ?)")
	res, err := ss.DbStore.Exec(insq, user.Email, user.FirstName, user.LastName, user.UserName, user.PassHash, user.PhotoURL)
	if err != nil {
		return nil, err// ErrInvalidInsert

	}

	//get the auto-assigned ID for the new row
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting new ID: %v\n", id)

	}

	user.ID = id

	return user, nil

}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ss *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	insq := "update users set first_name = ?, last_name = ? where id = ?"

	_, err := ss.DbStore.Exec(insq, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, ErrUserDNE

	}

	// get user after update?
	updatedUser, _ := ss.GetByID(id)

	return updatedUser, nil

}

func (ss *SQLStore) Delete(id int64) error {
	deleteString := "delete from users where id=?"
	_, err := ss.DbStore.Exec(deleteString, id)
	if err != nil {
		return fmt.Errorf("error deleting user with id %d: %v", id, err)

	}

	return nil

}

func (ss *SQLStore) InsertSignIn(id int64, currTime time.Time, ip string) error {
	insq := "insert into usersignin(id, whensignin, clientip) values (?, ?, ?)"
	_, err := ss.DbStore.Exec(insq, id, currTime, ip)
	if err != nil {
		return err

	}

	return nil

}
