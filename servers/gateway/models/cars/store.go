package cars

import (
	"errors"

)

// ErrNoCars returns when the user hasn't registered any cars yet
var ErrNoCars = errors.New("you have not registered any cars")
// ErrAlrRegist returns when the user tries to register a car that they already have before
var ErrAlrRegist = errors.New("you have already registered this car")
// ErrInvalidCar returns when the car requested is not registered to this user, or the car does not exist in the database
var ErrInvalidCar = errors.New("invalid car")

//Store represents a store for Cars
type Store interface {
	// InsertCar inserts the user's car
	InsertCar(newCar *Car, userid int64) (*Car, error)
	// GetCarByID returns the car with the given id
	GetCarByID(id int64) (*Car, error)
	// GetCarsByUserID returns a slice of all the cars they have registered
	GetCarsByUserID(userid int64) ([]*Car, error)
	// GetSpecificUserCar returns the car of the given id for the user
	GetSpecificUserCar(userid int64, carid int64) (*Car, error)
	// UpdateCar updates the registered vehicle's information
	UpdateCar(updates *UpdateCar, carid int64, userid int64) (*Car, error)
	// DeleteCar removes the registered car from the database for this particular user
	DeleteCarForUser(userid int64, carid int64) error

}
