package users

import (
	"errors"
	"time"
)

//ErrUserNotFound is returned when the user can't be found
var ErrUserNotFound = errors.New("user not found")

//ErrInvalidInsert is returned when the user can't be inserted
var ErrInvalidInsert = errors.New("cannot insert user")

//ErrUserDNE is returned when the user doesnt exist in dbms
var ErrUserDNE = errors.New("user not in database")

//Store represents a store for Users
type Store interface {
	//GetByID returns the User with the given ID
	GetByID(id int64) (*User, error)

	//GetByEmail returns the User with the given email
	GetByEmail(email string) (*User, error)

	//GetByUserName returns the User with the given Username
	GetByUserName(username string) (*User, error)

	//Insert inserts the user into the database, and returns
	//the newly-inserted User, complete with the DBMS-assigned ID
	Insert(user *User) (*User, error)

	//Update applies UserUpdates to the given user ID
	//and returns the newly-updated user
	Update(id int64, updates *Updates) (*User, error)

	//Delete deletes the user with the given ID
	Delete(id int64) error

	// InsertSignIn inserts every user log in attempt into the database
	// 	and returns any errors
	InsertSignIn(id int64, currTime time.Time, ip string) error
}
