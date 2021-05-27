package users

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID is a test function for the SQLStore's GetByID
func TestGetByID(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"No users in db yet",
			&User{},
			0,
			true,
		},
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &SQLStore{db}
		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"id",
			"email",
			"passhash",
			"username",
			"first_name",
			"last_name",
			"photourl"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,passhash,username,first_name,last_name,photourl from users where id=?"
		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else { // expected an error
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)
			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}

			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByEmail is a test function for the SQLStore's GetByEmail
func TestGetByEmail(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet   string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"test@test.com",
			false,
		},
		{
			"User Not Found",
			&User{},
			"notfound@gmail.com",
			true,
		},
		{
			"User with longer email found",
			&User{
				1234567890,
				"anothervalidemail@com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"anothervalidemail@com",
			false,
		},
		{
			"User with very similar email not found",
			&User{
				3,
				"tesl@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"someurl",
			},
			"tesl@test.com",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &SQLStore{db}
		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"id",
			"email",
			"passhash",
			"username",
			"first_name",
			"last_name",
			"photourl"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// TODO: update to match the query used in your Store implementation
		query := "select id,email,passhash,username,first_name,last_name,photourl from users where email=?"
		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else { // expected an error
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}

			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByUserName is a test function for the SQLStore's GetByUserName
func TestGetByUserName(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		unToGet      string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"username",
			false,
		},
		{
			"User Not Found",
			&User{},
			"someuser",
			true,
		},
		{
			"User with longer email found",
			&User{
				1234567890,
				"anothervalidemail@com",
				[]byte("passhash123"),
				"username229908019240192",
				"firstname",
				"lastname",
				"photourl",
			},
			"username229908019240192",
			false,
		},
		{
			"User with very similar email not found",
			&User{
				3,
				"tesl@test.com",
				[]byte("passhash123"),
				"username229908019240194",
				"firstname",
				"lastname",
				"someurl",
			},
			"username229908019240193",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &SQLStore{db}
		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"id",
			"email",
			"passhash",
			"username",
			"first_name",
			"last_name",
			"photourl"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,passhash,username,first_name,last_name,photourl from users where username=?"
		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.unToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByUserName(c.unToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
			
		} else { // expected an error
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.unToGet).WillReturnRows(row)
			user, err := mainSQLStore.GetByUserName(c.unToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}

			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestInsert is a test function for the SQLStore's Insert
func TestInsert(t *testing.T) { // test injection attacks
	// Create a slice of test cases
	cases := []struct {
		name         string
		passedUser   *User
		expectedUser *User
		expectUN     string
		expectError  bool
	}{
		{
			"Insert proper user",
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			&User{
				ID:        1,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			"username",
			false,
		},
		{
			"Insert another proper user",
			&User{
				Email:     "what@test.com",
				PassHash:  []byte("yeyeyeye"),
				UserName:  "someuser",
				FirstName: "aname",
				LastName:  "anothername",
				PhotoURL:  "anotherphotourl",
			},
			&User{
				ID:        1,
				Email:     "what@test.com",
				PassHash:  []byte("yeyeyeye"),
				UserName:  "someuser",
				FirstName: "aname",
				LastName:  "anothername",
				PhotoURL:  "anotherphotourl",
			},
			"someuser",
			false,
		},
		{
			"Insert bad injection attack attempt",
			&User{
				Email:     "defbad@info.com",
				PassHash:  []byte("some secret password"),
				UserName:  "mal",
				FirstName: "; drop table users; select ''",
				LastName:  "; drop table users; select ''",
				PhotoURL:  "; drop table users; select ''",
			},
			&User{
				ID:        2,
				Email:     "defbad@info.com",
				PassHash:  []byte("some secret password"),
				UserName:  "mal",
				FirstName: "; drop table users; select ''",
				LastName:  "; drop table users; select ''",
				PhotoURL:  "; drop table users; select ''",
			},
			"mal",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)

		}

		defer db.Close()

		mainSQLStore := &SQLStore{db}

		query := "insert into users(email, first_name, last_name, username, passhash, photourl) values (?, ?, ?, ?, ?, ?)"
		if c.expectError {
			mock.ExpectExec(query).WithArgs(c.expectedUser.Email, c.expectedUser.FirstName, c.expectedUser.LastName, c.expectedUser.UserName, c.expectedUser.PassHash, c.expectedUser.PhotoURL).WillReturnError(ErrInvalidInsert)

			dbmsUser, err := mainSQLStore.Insert(c.passedUser)
			if dbmsUser != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrInvalidInsert, err)

			}

		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectExec(query).WithArgs(c.expectedUser.Email, c.expectedUser.FirstName, c.expectedUser.LastName, c.expectedUser.UserName, c.expectedUser.PassHash, c.expectedUser.PhotoURL).WillReturnResult(sqlmock.NewResult(c.expectedUser.ID, 1)) //(c.expectedUser.ID, 1) //WillReturnRows(row)

			dbmsUser, err := mainSQLStore.Insert(c.passedUser) // inserts into the database
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)

			}

			if !reflect.DeepEqual(dbmsUser, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)

			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)

		}
	}

}

// TestInsertAndDelete is a test function for the SQLStore's Insert & Delete
func TestInsertAndDelete(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		passedUser   *User
		passedID     int64
		idToDel      int64
		expectDelErr bool
	}{
		{
			"Insert proper user to successfully delete 1",
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			1,
			1,
			false,
		},
		{
			"Insert proper user but delete user not in database 1",
			&User{
				Email:     "test2@test.com",
				PassHash:  []byte("woooooooo3"),
				UserName:  "someuser",
				FirstName: "yeah",
				LastName:  "whoooo",
				PhotoURL:  "someurl.com/yeah",
			},
			2,
			4,
			true,
		},
		{
			"Insert proper user but delete user not in database 2",
			&User{
				Email:     "somegal@test.com",
				PassHash:  []byte("supersecretpasshash938042937298"),
				UserName:  "yeah",
				FirstName: "okay",
				LastName:  "youdont",
				PhotoURL:  "gravatar/weoirjw",
			},
			3,
			20319480914032,
			true,
		},
		{
			"Insert proper user to successfully delete 2",
			&User{
				Email:     "coolio@test.com",
				PassHash:  []byte("itssocold"),
				UserName:  "wowee",
				FirstName: "washington",
				LastName:  "george",
				PhotoURL:  "somefancyurl",
			},
			4,
			4,
			false,
		},
		{
			"Insert then delete previously added one",
			&User{
				Email:     "somegal@test.com",
				PassHash:  []byte("supersecretpasshash938042937298"),
				UserName:  "yeah",
				FirstName: "okay",
				LastName:  "youdont",
				PhotoURL:  "gravatar/weoirjw",
			},
			5,
			2,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)

		}

		defer db.Close()

		mainSQLStore := &SQLStore{db}

		query := "insert into users(email, first_name, last_name, username, passhash, photourl) values (?, ?, ?, ?, ?, ?)"
		mock.ExpectExec(query).WithArgs(c.passedUser.Email, c.passedUser.FirstName, c.passedUser.LastName, c.passedUser.UserName, c.passedUser.PassHash, c.passedUser.PhotoURL).WillReturnResult(sqlmock.NewResult(c.passedID, 1))

		_, err = mainSQLStore.Insert(c.passedUser)
		if err != nil {
			t.Errorf("Unexpected error [%v]", err)

		}

		query = "delete from users where id=?"

		if !c.expectDelErr {
			mock.ExpectExec(query).WithArgs(c.idToDel).WillReturnResult(sqlmock.NewResult(c.idToDel, 1))

			err = mainSQLStore.Delete(c.idToDel)
			if err != nil {
				t.Errorf("Unexpected error on test [%s]: %v", c.name, err)

			}

		} else { // we do expect error
			mock.ExpectExec(query).WithArgs(c.idToDel).WillReturnError(ErrUserDNE)

			err = mainSQLStore.Delete(c.idToDel) // inserts into the database
			if err == nil {
				t.Errorf("Expected error [%v] on test [%s] but didn't get one", ErrUserDNE, c.name)

			}

		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)

		}
	}

}

// TestUpdate is a test function for the SQLStore's Update
func TestUpdate(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		userUpdates  *Updates
		idToUpd      int64
		expectedUser *User
		expectError bool
	}{
		{
			"Proper update",
			&Updates{
				FirstName: "Sarah",
				LastName: "Gefiso",
			},
			1,
			&User{
				ID: 1,
				Email: "test@test.com",
				PassHash: []byte("passhash123"),
				UserName: "username",
				FirstName: "Sarah",
				LastName: "Gefiso",
				PhotoURL: "photourl",
			},
			false,
		},
		{
			"Proper update 2",
			&Updates{
				FirstName: "NoLastName",
				LastName: "",
			},
			2,
			&User{
				ID: 2,
				Email: "sosad@test.com",
				PassHash: []byte("poeijgapoieh"),
				UserName: "Oc",
				FirstName: "NoLastName",
				LastName: "",
				PhotoURL: "photourl",
			},
			false,
		},
		{
			"Updating User that Doesnt exist",
			&Updates{
				FirstName: "saoiefj",
				LastName: "asjvasdf",
			},
			10,
			&User{},
			true,
		},

	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}

		defer db.Close()

		mainSQLStore := &SQLStore{db}
		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"id",
			"email",
			"passhash",
			"username",
			"first_name",
			"last_name",
			"photourl"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		//Matching the query used in Store implementation
		query := "update users set first_name = ?, last_name = ? where id = ?"
	  query2 := "select id,email,passhash,username,first_name,last_name,photourl from users where id=?"

		if c.expectError {
			mock.ExpectExec(query).WithArgs(c.userUpdates.FirstName, c.userUpdates.LastName, c.idToUpd).WillReturnError(ErrUserDNE)

			_, err := mainSQLStore.Update(c.idToUpd, c.userUpdates)
			if err == nil {
				t.Errorf("Expected error [%v] on test [%s] but didn't get one", ErrUserDNE, c.name)

			}

		} else { // expected an error
			mock.ExpectExec(query).WithArgs(c.userUpdates.FirstName, c.userUpdates.LastName, c.idToUpd).WillReturnResult(sqlmock.NewResult(c.idToUpd, 1))
			mock.ExpectQuery(query2).WithArgs(c.idToUpd).WillReturnRows(row)

			newUser, err := mainSQLStore.Update(c.idToUpd, c.userUpdates)
			if err != nil {
				t.Errorf("Unexpected error when updating: %v", err)
			}
	
			if !reflect.DeepEqual(newUser, c.expectedUser) {
				fmt.Println(newUser)
				fmt.Println(c.expectedUser)
				t.Errorf("Error, invalid match in test [%s]", c.name)
	
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}