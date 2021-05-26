package handlers

import (
	"bytes"
	"fmt"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/choijos/assignments-choijos/servers/gateway/models/users"
	"github.com/choijos/assignments-choijos/servers/gateway/sessions"

)

// TestMySQLStore is used for tests
type TestMySQLStore struct {

}

// NewFakeStore creates a new TestMySQLStore
func NewFakeStore() *TestMySQLStore {
	return &TestMySQLStore{}

}

// GetByID is a fake function for TestMySQLStore
func (s *TestMySQLStore) GetByID(id int64) (*users.User, error) {
	if id == 1 {
		return &users.User{
			ID: 1,
		}, nil

	}

	return nil, errors.New("only id 1 can be found")

}

// GetByEmail is a fake function for TestMySQLStore
func (s *TestMySQLStore) GetByEmail(email string) (*users.User, error) {
	if email == "emailcannotfind@email.com" {
		return nil, users.ErrUserNotFound

	}

	if email == "emailfailedcreds@email.com" {
		return &users.User{ID: 1, PassHash: []byte("not the right passcode")}, errors.New("error getting user with given email and credentials")

	}

	if email == "emailworkingcreds@email.com" {
		user := &users.User{ID: 1}
		user.SetPassword("password")
		return user, nil

	}

	return &users.User{
		ID: 1,
	}, nil

}

// GetByUserName is a fake function for TestMySQLStore
func (s *TestMySQLStore) GetByUserName(username string) (*users.User, error) {
	return &users.User{
		ID: 1,
	}, nil

}

// Insert is a fake function for TestMySQLStore
func (s *TestMySQLStore) Insert(user *users.User) (*users.User, error) {
	return &users.User{
		ID: 1,
	}, nil

}

// InsertSignIn is a fake function for TestMySQLStore
func (s *TestMySQLStore) InsertSignIn(id int64, currTime time.Time, ip string) error {
	return nil

}

// Update is a fake function for TestMySQLStore
func (s *TestMySQLStore) Update(id int64, updates *users.Updates) (*users.User, error) {
	return &users.User{
		ID: 1,
	}, nil

}

// Delete is a fake function for TestMySQLStore
func (s *TestMySQLStore) Delete(id int64) error { return nil }

// TestUsersHandler performs tests on the UsersHandler handler
func TestUsersHandler(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		reqBody     []byte
		ctype       string
		fakesignkey string
		errType     int
		expectError bool
	}{
		{
			"Invalid method - GET",
			"GET",
			nil,
			"",
			"asdf",
			405,
			true,
		},
		{
			"Invalid method - PATCH",
			"PATCH",
			nil,
			"",
			"asdf",
			405,
			true,
		},
		{
			"Invalid content type",
			"POST",
			nil,
			"",
			"asdf",
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"Request Body not correct format",
			"POST",
			[]byte(`something that isnt json, so it should not be read/decoded properly`),
			"application/json",
			"asdf",
			http.StatusInternalServerError,
			true,
		},
		{
			"Request Body not correct format",
			"POST",
			[]byte(`{something that isnt json, so it should not be read/decoded properly}`),
			"application/json",
			"asdf",
			http.StatusInternalServerError,
			true,
		},
		{
			"Invalid new user results in user conversion error (ToUser)",
			"POST",
			[]byte(`{"email": "test@gmail.com", "password": "badpassword", "passwordConf": "badpasswordconf", "userName": "GoodGuy"}`),
			"application/json",
			"asdf",
			http.StatusBadRequest,
			true,
		},
		{
			"Empty session key",
			"POST",
			[]byte(`{"email": "test@gmail.com", "password": "password", "passwordConf": "password", "userName": "GoodGuy"}`),
			"application/json",
			"",
			http.StatusInternalServerError,
			true,
		},
		{
			"Valid/correct case",
			"POST",
			[]byte(`{"email": "test@gmail.com", "password": "password", "passwordConf": "password", "userName": "GoodGuy"}`),
			"application/json",
			"asdf",
			http.StatusCreated,
			false,
		},
	}

	for _, c := range cases {
		memstore := sessions.NewMemStore(time.Hour, time.Minute)
		userstore := NewFakeStore()

		newCtx := &HandlerContext{
			c.fakesignkey,
			memstore,
			userstore,
		}

		req, err := http.NewRequest(c.method, "/v1/users", bytes.NewReader(c.reqBody))
		if err != nil {
			t.Errorf("error trying to create new request for test %s: %v", c.name, err)

		}

		req.Header.Set("Content-Type", c.ctype)
		rr := httptest.NewRecorder()
		newCtx.UsersHandler(rr, req)

		if c.expectError {
			if rr.Code != c.errType {
				t.Errorf("Incorrect http code returned for error case %s\nExpected: %d\nActual: %d", c.name, c.errType, rr.Code)

			}

		} else {
			if rr.Code != http.StatusCreated {
				t.Errorf("Incorrect http code returned for correct case %s\nExpected: %d\nActual: %d", c.name, http.StatusCreated, rr.Code)

			}

		}

	}

}

// TestSpecificUserHandler performs tests on the SpecificUserHandler handler
func TestSpecificUserHandler(t *testing.T) {
	cases := []struct {
		name    string
		method  string
		urlBase string
		ctype       string
		fakesignkey string
		reqBody     []byte
		errType     int
		expectError bool
	}{
		{
			"Unsupported Method",
			"POST",
			"1",
			"",
			"asdf",
			nil,
			http.StatusMethodNotAllowed,
			true,
		},
		{
			"GET - User not authorized",
			"GET",
			"",
			"",
			"asdf",
			nil,
			http.StatusUnauthorized,
			true,
		},
		{
			"GET - Invalid string base url",
			"GET",
			"notgood",
			"",
			"asdf",
			nil,
			http.StatusNotAcceptable,
			true,
		},
		{
			"GET - Cannot get requested user of id 0",
			"GET",
			"0",
			"",
			"asdf",
			nil,
			http.StatusNotFound,
			true,
		},
		{
			"GET - Successful user handling with a user id base url",
			"GET",
			"1",
			"",
			"asdf",
			nil,
			http.StatusOK,
			false,
		},
		{
			"GET - Successful user handling with 'me' a  base url",
			"GET",
			"me",
			"",
			"asdf",
			nil,
			http.StatusOK,
			false,
		},
		{
			"PATCH - URL not 'me' or currently authenticated user",
			"PATCH",
			"5",
			"",
			"asdf",
			nil,
			http.StatusForbidden,
			true,
		},
		{
			"PATCH - Unsupported content type",
			"PATCH",
			"1",
			"notjson",
			"asdf",
			nil,
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"PATCH - Bad nil request body",
			"PATCH",
			"1",
			"application/json",
			"asdf",
			nil,
			http.StatusBadRequest,
			true,
		},
		{
			"PATCH - Incorrect Fields to update",
			"PATCH",
			"1",
			"application/json",
			"asdf",
			[]byte(`{"firstName": "just2"}`),
			http.StatusInternalServerError,
			true,
		},
		{
			"PATCH - Successful case",
			"PATCH",
			"me",
			"application/json",
			"asdf",
			[]byte(`{"firstName": "just"}`),
			http.StatusOK,
			false,
		},
		
	}

	for _, c := range cases {
		memstore := sessions.NewMemStore(time.Hour, time.Minute)
		userstore := NewFakeStore()

		newCtx := &HandlerContext{
			c.fakesignkey,
			memstore,
			userstore,
		}

		req, err := http.NewRequest(c.method, fmt.Sprintf("/v1/users/%s", c.urlBase), nil)

		if c.method == "PATCH" {
			req, err = http.NewRequest(c.method, fmt.Sprintf("/v1/users/%s", c.urlBase), bytes.NewReader(c.reqBody))

		}

		if err != nil {
			t.Errorf("error trying to create new request for test %s: %v", c.name, err)

		}

		if c.method == "PATCH" {
			req.Header.Set("Content-Type", c.ctype)

		}

		rr := httptest.NewRecorder()
		fakeUser, err := userstore.GetByID(1)

		if err != nil {
			t.Errorf("Unexpected error when getting user from fake store case %s: %v", c.name, err)

		}

		newSession := SessionState{time.Now(), fakeUser}
		
		if c.name != "GET - User not authorized" {
			sessID, err := sessions.BeginSession(c.fakesignkey, memstore, newSession, rr)
			if err != nil {
				t.Errorf("Unwanted error when beginning session for case %s: %v", c.name, err)

			} else {
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sessID))

			}

		}

		newCtx.SpecificUserHandler(rr, req)

		if c.expectError {
			if rr.Code != c.errType {
				t.Errorf("Incorrect http code returned for error case %s\nExpected: %d\nActual: %d", c.name, c.errType, rr.Code)

			}

		} else {
			if rr.Code != c.errType {
				t.Errorf("Incorrect http code returned for correct case %s\nExpected: %d\nActual: %d", c.name, c.errType, rr.Code)

			} else if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Content type header not set to json for case %s\nActual: %s", c.name, rr.Header().Get("Content-Type"))

			}

		}

	}

}

// TestSessionsHandler performs tests on the SessionsHadler handler
func TestSessionsHandler(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		reqBody     []byte
		ctype       string
		fakesignkey string
		errType     int
		expectError bool
	}{
		{
			"Invalid method",
			"GET",
			nil,
			"",
			"asdf",
			http.StatusMethodNotAllowed,
			true,
		},
		{
			"Incorrect content type",
			"POST",
			nil,
			"notjson",
			"asdf",
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"Error decoding nil request body",
			"POST",
			nil,
			"application/json",
			"asdf",
			http.StatusBadRequest,
			true,
		},
		{
			"Error getting user by email - does not exist in store",
			"POST",
			[]byte(`{"email": "emailcannotfind@email.com"}`),
			"application/json",
			"asdf",
			http.StatusUnauthorized,
			true,
		},
		{
			"Error getting user by email",
			"POST",
			[]byte(`{"email": "emailfailedcreds@email.com"}`),
			"application/json",
			"asdf",
			http.StatusInternalServerError,
			true,
		},
		{
			"Empty signing key, cannot begin session",
			"POST",
			[]byte(`{"email": "emailworkingcreds@email.com"}`),
			"application/json",
			"",
			http.StatusInternalServerError,
			true,
		},
		{
			"Successful case",
			"POST",
			[]byte(`{"email": "emailworkingcreds@email.com"}`),
			"application/json",
			"asdf",
			http.StatusCreated,
			false,
		},
		
	}

	for _, c := range cases {
		memstore := sessions.NewMemStore(time.Hour, time.Minute)
		userstore := NewFakeStore()

		newCtx := &HandlerContext{
			c.fakesignkey,
			memstore,
			userstore,
		}

		req, err := http.NewRequest(c.method, "/v1/sessions", bytes.NewReader(c.reqBody))
		if err != nil {
			t.Errorf("error trying to create new request for test %s: %v", c.name, err)

		}

		if c.name == "Successful case" {
			req.Header.Set("X-Forwarded-For", "9.2312.3, 204.212")

		}

		req.Header.Set("Content-Type", c.ctype)
		rr := httptest.NewRecorder()
		newCtx.SessionsHandler(rr, req)

		if c.expectError {
			if rr.Code != c.errType {
				t.Errorf("Incorrect http code returned for error case %s\nExpected: %d\nActual: %d", c.name, c.errType, rr.Code)

			}

		} else {
			if rr.Code != http.StatusCreated {
				t.Errorf("Incorrect http code returned for correct case %s\nExpected: %d\nActual: %d", c.name, http.StatusCreated, rr.Code)

			}

		}

	}

}

// TestSpecificSessionHandler performs tests on the SpecificSessionHandler handler
func TestSpecificSessionHandler(t *testing.T) {
	cases := []struct {
		name    string
		method  string
		urlBase string
		fakesignkey string
		errType     int
		expectError bool
	}{
		{
			"Unsupported Method",
			"POST",
			"",
			"asdf",
			http.StatusMethodNotAllowed,
			true,
		},
		{
			"Invalid base url path",
			"DELETE",
			"notmine",
			"asdf",
			http.StatusForbidden,
			true,
		},
		{
			"Error ending session - no header auth set",
			"DELETE",
			"mine",
			"asdf",
			http.StatusInternalServerError,
			true,
		},
		{
			"Successful case",
			"DELETE",
			"mine",
			"asdf",
			0,
			false,
		},
		
	}

	for _, c := range cases {
		memstore := sessions.NewMemStore(time.Hour, time.Minute)
		userstore := NewFakeStore()

		newCtx := &HandlerContext{
			c.fakesignkey,
			memstore,
			userstore,
		}

		req, err := http.NewRequest(c.method, fmt.Sprintf("/v1/sessions/%s", c.urlBase), nil)

		if err != nil {
			t.Errorf("error trying to create new request for test %s: %v", c.name, err)

		}

		rr := httptest.NewRecorder()
		fakeUser, err := userstore.GetByID(1)

		if err != nil {
			t.Errorf("Unexpected error when getting user from fake store case %s: %v", c.name, err)

		}

		newSession := SessionState{time.Now(), fakeUser}

		sessID, err := sessions.BeginSession(c.fakesignkey, memstore, newSession, rr)
		if err != nil {
			t.Errorf("Unwanted error when beginning session for case %s: %v", c.name, err)

		}

		if c.name != "Error ending session - no header auth set" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sessID))

		}

		newCtx.SpecificSessionHandler(rr, req)

		if c.expectError {
			if rr.Code != c.errType {
				t.Errorf("Incorrect http code returned for error case %s\nExpected: %d\nActual: %d", c.name, c.errType, rr.Code)

			}

		}

	}

}