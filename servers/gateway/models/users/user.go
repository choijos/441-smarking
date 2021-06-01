package users

import (
	"crypto/md5"
	"fmt"
	"net/mail"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID          int64  `json:"id"`
	Email       string `json:"-"` //never JSON encoded/decoded
	PassHash    []byte `json:"-"` //never JSON encoded/decoded
	UserName    string `json:"userName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhotoURL    string `json:"photoURL"`
	PhoneNumber string `json:"phoneNumber"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	PhoneNumber  string `json:"phoneNumber"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("email not valid")

	}

	_, err = strconv.ParseInt(nu.PhoneNumber, 10, 64)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("error converting provided User ID from url to int64: %v", err), http.StatusNotAcceptable)
	// 	return

	// }

	if len(nu.PhoneNumber) < 10 || err != nil {
		return fmt.Errorf("please enter a valid phone number")

	}

	if len(nu.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")

	}

	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("passwords do not match")

	}

	checkSpaces := strings.Split(nu.UserName, " ")

	if cap(checkSpaces) != 1 {
		return fmt.Errorf("username must use no spaces")

	}

	if len(nu.UserName) == 0 {
		return fmt.Errorf("username must be non-zero length")

	}

	return nil
}

// ToUser converts the NewUser to a User, setting the
// PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	err := nu.Validate()
	if err != nil {
		return nil, err

	}

	phone := "+1" + nu.PhoneNumber

	newUserStruct := &User{
		Email:       nu.Email,
		UserName:    nu.UserName,
		FirstName:   nu.FirstName,
		LastName:    nu.LastName,
		PhoneNumber: phone,
	}

	newUserStruct.PhotoURL = fmt.Sprintf("%s%s", gravatarBasePhotoURL, fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(strings.TrimSpace(newUserStruct.Email))))))
	newUserStruct.SetPassword(nu.Password)

	return newUserStruct, nil

}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	if len(u.FirstName) != 0 && len(u.LastName) != 0 {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)

	} else if len(u.FirstName) != 0 {
		return u.FirstName

	} else if len(u.LastName) != 0 {
		return u.LastName

	}

	return ""

}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")

	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error generating bcrypt hash: %v\n", err)

	}

	u.PassHash = hash

	return nil

}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password)); err != nil {
		return fmt.Errorf("password doesn't match stored hash")

	}

	return nil

}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	if updates == nil {
		return fmt.Errorf("no updates made")

	}

	_, err := strconv.ParseInt(updates.PhoneNumber, 10, 64)

	if len(updates.PhoneNumber < 10) || err != nil {
		return fmt.Errorf("enter a valid phone number")

	}

	if strings.ContainsAny(updates.FirstName, "1234567890") {
		return fmt.Errorf("first name cannot contain any numbers")

	} else if strings.ContainsAny(updates.LastName, "1234567890") {
		return fmt.Errorf("last name cannot contain any numbers")

	}

	u.PhoneNumber = "+1" + updates.PhoneNumber
	u.FirstName = updates.FirstName
	u.LastName = updates.LastName

	return nil

}
