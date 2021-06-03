package cars

import (
	"database/sql"
	"fmt"
)

// SQLStore represents a mysql database containing all cars
type SQLStore struct { // probably need to figure out how to import this from the users mysqlstore.go
	DbStore *sql.DB
}

// For this function, we might have to pull the userid from the currently enticated user rather than passing it in as a parameter?
//
// InsertCar adds the given car to the database. Returns the car struct
// 	with its new DBMS assigned ID. Returns error if the user has already registered
// 	this car before
func (ss *SQLStore) InsertCar(newCar *Car, userid int64) (*Car, error) {
	// might have to check before the actual executed query if the licenseplate is already registered with this particular user?
	rows, err := ss.DbStore.Query("select ID, LicensePlate from cars where UserID = ? and LicensePlate = ?", userid, newCar.LicensePlate)
	if err != nil {
		return nil, err

	}

	testCar := Car{}

	for rows.Next() {
		err := rows.Scan(&testCar.ID)
		if err != nil {
			return nil, ErrAlrRegist

		}

	}

	if (Car{}) != testCar {
		return nil, ErrAlrRegist

	}

	ins := "insert into cars(LicensePlate, UserID, Make, Model, ModelYear, Color) values (?, ?, ?, ?, ?, ?)"
	res, err := ss.DbStore.Exec(ins, newCar.LicensePlate, userid, newCar.Make, newCar.Model, newCar.ModelYear, newCar.Color)
	if err != nil {
		return nil, err

	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting new car ID: %v", err)

	}

	newCar.ID = id
	newCar.UserID = userid

	return newCar, nil

}

// dont know if we actually will need this function
//
// GetCarByID returns a Car struct containing the database information on the Car with
// 	the given ID
func (ss *SQLStore) GetCarByID(id int64) (*Car, error) {
	rows, err := ss.DbStore.Query("select ID, LicensePlate, UserID, Make, Model, ModelYear, Color from cars where ID = ?", id)
	if err != nil {
		return nil, fmt.Errorf("unable to grab car of given id: %v", err)

	}

	// testCar := Car{}

	// for rows.Next() {
	// 	err = rows.Scan(&testCar.ID)
	// 	if err == sql.ErrNoRows {
	// 		return nil, ErrAlrRegist

	// 	} else if err != nil {
	// 		return nil, err

	// 	}

	// }

	defer rows.Close()

	retCar := Car{}

	for rows.Next() {
		if err := rows.Scan(&retCar.ID, &retCar.LicensePlate, &retCar.UserID, &retCar.Make, &retCar.Model, &retCar.ModelYear, &retCar.Color); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)

		}

	}

	if (Car{}) == retCar {
		return nil, ErrInvalidCar

	}

	return &retCar, nil

}

// probably should return an array of structs
//
// GetCarsByUserID returns a slice of Car structs containing the information on all
// 	of the cars this user with the given ID has registered
func (ss *SQLStore) GetCarsByUserID(userid int64) ([]*Car, error) {
	rows, err := ss.DbStore.Query("select ID, LicensePlate, UserID, Make, Model, ModelYear, Color from cars where UserID = ?", userid)
	if err != nil {
		return nil, ErrNoCars

	}

	defer rows.Close()

	allCars := []*Car{}
	for rows.Next() {
		oneCar := Car{}
		// if err := rows.Scan(); err == sql.ErrNoRows {
		// 	return nil, ErrInvalidCar

		// }
		err := rows.Scan(&oneCar.ID, &oneCar.LicensePlate, &oneCar.UserID, &oneCar.Make, &oneCar.Model, &oneCar.ModelYear, &oneCar.Color)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)

		} 

		allCars = append(allCars, &oneCar)

	}

	// if (Car{}) == theCar {
	// 	return nil, ErrInvalidCar

	// }

	return allCars, nil

}

// Might not want to return all columns?
//
// GetSpecificUserCar returns a Car struct containing the database information for the
// 	car with the given ID for the user
func (ss *SQLStore) GetSpecificUserCar(userid int64, carid int64) (*Car, error) {
	rows, err := ss.DbStore.Query("select ID, LicensePlate, UserID, Make, Model, ModelYear, Color from cars where ID = ? and UserID = ?", carid, userid)
	if err != nil {
		return nil, ErrNoCars

	}

	defer rows.Close()

	theCar := Car{}
	for rows.Next() {
		err := rows.Scan(&theCar.ID, &theCar.LicensePlate, &theCar.UserID, &theCar.Make, &theCar.Model, &theCar.ModelYear, &theCar.Color)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)

		}

	}

	if (Car{}) == theCar {
		return nil, ErrInvalidCar

	}

	return &theCar, nil

}

// UpdateCar applies the passed in updates for this car registered under the user in the database,
// 	returns the updated car's information in a Car struct
func (ss *SQLStore) UpdateCar(updates *UpdateCar, carid int64, userid int64) (*Car, error) {
	// _, err := ss.GetSpecificUserCar(userid, carid)
	// if err != nil {
	// 	return nil, err

	// }

	if updates == nil {
		return nil, fmt.Errorf("no updates made")

	}

	if updates.LicensePlate != "" {
		ins := "update cars set LicensePlate = ? where ID = ? and UserID = ?"
		_, err := ss.DbStore.Exec(ins, updates.LicensePlate, carid, userid)
		if err != nil {
			return nil, ErrNoCars // check
	
		}

	}

	if updates.Make != "" {
		ins := "update cars set Make = ? where ID = ? and UserID = ?"
		_, err := ss.DbStore.Exec(ins, updates.Make, carid, userid)
		if err != nil {
			return nil, ErrNoCars // check
	
		}

	}

	if updates.Model != "" {
		ins := "update cars set Model = ? where ID = ? and UserID = ?"
		_, err := ss.DbStore.Exec(ins, updates.Model, carid, userid)
		if err != nil {
			return nil, ErrNoCars // check
	
		}

	}

	if updates.ModelYear != "" {
		ins := "update cars set ModelYear = ? where ID = ? and UserID = ?"
		_, err := ss.DbStore.Exec(ins, updates.ModelYear, carid, userid)
		if err != nil {
			return nil, ErrNoCars // check
	
		}

	}

	if updates.Color != "" {
		ins := "update cars set Color = ? where ID = ? and UserID = ?"
		_, err := ss.DbStore.Exec(ins, updates.Color, carid, userid)
		if err != nil {
			return nil, ErrNoCars // check
	
		}

	}

	updatedCar, err := ss.GetCarByID(carid)
	if err != nil {
		return nil, err

	}

	return updatedCar, nil

}

// DeleteCarForUser removes the car with given ID for the user from the database
func (ss *SQLStore) DeleteCarForUser(userid int64, carid int64) error {
	del := "delete from cars where ID = ? and UserID = ?"
	_, err := ss.DbStore.Exec(del, carid, userid)
	if err != nil {
		return ErrInvalidCar

	}

	return nil

}