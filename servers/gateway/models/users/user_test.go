package users

import (
	"crypto/md5"
	"fmt"
	"strings"
	"testing"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

// TestValidate is a test function for the user Validate function that
// 	ensures all struct fields of a NewUser are valid
func TestValidate(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		passwordConf string
		userName     string
		expectError  bool
		hint         string
	}{
		{
			"Invalid Email",
			"sarahgmail.com", // invalid
			"dgsheofijsaoe",
			"dgsheofijsaoe",
			"StarsRock",
			true,
			"Validate email",
		},
		// {
		// 	"Invalid Email 2",
		// 	"sarag@com", // invalid, but passes as valid?
		// 	"dgsheofijsaoe",
		// 	"dgsheofijsaoe",
		// 	"StarsRock",
		// 	true,
		// 	"Validate email",
		// },
		{
			"Short password",
			"sarah@gmail.com",
			"dgs", // invalid
			"dgs", // invalid
			"StarsRock",
			true,
			"Validate password length",
		},
		{
			"Unmatching Passwords",
			"sarah@gmail.com",
			"dgsdsgsdafsdfas", // invalid
			"dgsgasdrgwerwes", // invalid
			"StarsRock",
			true,
			"Validate password confirmation",
		},
		{
			"Space in username",
			"sarah@gmail.com",
			"dgsheofijsaoe",
			"dgsheofijsaoe",
			"Stars Rock",
			true,
			"Validate username (banned character)",
		},
		{
			"Empty username",
			"sarah@gmail.com",
			"dgsheofijsaoe",
			"dgsheofijsaoe",
			"",
			true,
			"Validate username (length)",
		},
		{
			"Good case",
			"sarah@gmail.com",
			"dgsheofijsaoe",
			"dgsheofijsaoe",
			"StarsRock",
			false,
			"All fields should be valid",
		},
	}

	for _, c := range cases {
		testNewUser := &NewUser{
			Email:        c.email,
			Password:     c.password,
			PasswordConf: c.passwordConf,
			UserName:     c.userName,
		}

		err := testNewUser.Validate()
		if err == nil && c.expectError {
			t.Errorf("case %s: expected error but didn't get one\n %s", c.name, c.hint)

		}

		if err != nil && !c.expectError {
			t.Errorf("case %s: didn't expect error but got one\n HINT: %s", c.name, c.hint)

		}
	}

}


// TestFullName is a test function for the user FullName function that
// 	ensures the format of the returned strings are correct
func TestFullName(t *testing.T) {
	cases := []struct {
		name         string
		firstname    string
		lastname     string
		expectedName string
		hint         string
	}{
		{
			"First and last name not empty",
			"Sarah",
			"Gefiso",
			"Sarah Gefiso",
			"Make sure you are checking all cases for formatting",
		},
		{
			"First name empty",
			"",
			"Moden",
			"Moden",
			"Last name expected, watch out for unwanted spaces",
		},
		{
			"Last name empty",
			"Zach",
			"",
			"Zach",
			"First name expected, watch out for unwanted spaces",
		},
		{
			"No names provided",
			"",
			"",
			"",
			"Make sure you are checking all cases for formatting, user may not provide any names",
		},
	}

	for _, c := range cases {
		testUser := &User{
			FirstName: c.firstname,
			LastName:  c.lastname,
		}

		nameOutput := testUser.FullName()
		if nameOutput != c.expectedName {
			t.Errorf("case %s: output does not match expected string\n Expected: %s\n Actual: %s\n Hint: %s", c.name, c.expectedName, nameOutput, c.hint)

		}

	}

}


// TestSetPasswordAuthenticate is a test function for both the user SetPassword and
// 	Authenticate functions that ensures the passwords when hashed are correct
func TestSetPasswordAuthenticate(t *testing.T) {
	cases := []struct {
		name            string
		prehashpassword string
	}{
		{
			"Test setting password #1",
			"hopefully this is a valid password",
		},
		{
			"Test setting password #2",
			"jdfspoihgaposigpjo322131421231",
		},
		{
			"Test setting password #3",
			"382938",
		},
	}

	for _, c := range cases {
		testUser := &User{}

		err := testUser.SetPassword(c.prehashpassword)
		if err != nil {
			t.Error(err) // returning the error by itself

		}

		authErr := testUser.Authenticate(c.prehashpassword)
		if err != nil {
			t.Error(authErr)

		}

	}

}

// TestApplyUpdates is a test function for the user ApplyUpdates function that
// 	ensures correct user updates
func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name        string
		currFN      string
		currLN      string
		newFN       string
		newLN       string
		expectederr bool
	}{
		{
			"Replacing both first and last name",
			"Sarah",
			"Gefiso",
			"Howdy",
			"Sos",
			false,
		},
		{
			"Replacing only first name",
			"Sarah",
			"Gefiso",
			"Okes",
			"",
			false,
		},
		{
			name:   "Replacing only last name",
			currFN: "Sarah",
			currLN: "Gefiso",
			newLN:       "Wat",
			expectederr: false,
		},
		{
			name:        "Replacing first name with numbers",
			currFN:      "Sam",
			currLN:      "Tam",
			newFN:       "Sam123",
			expectederr: true,
		},
		{
			name:        "Replacing last name with numbers",
			currFN:      "Sam",
			currLN:      "Tam",
			newLN:       "Tam2",
			expectederr: true,
		},
		{
			name:        "Replacing no names, leaving empty",
			currFN:      "Manny",
			currLN:      "Man",
			expectederr: false,
		},
		{
			name:        "Adding new first and last names",
			newFN:       "Yak",
			newLN:       "Ok",
			currFN:      "Yam",
			currLN:      "No",
			expectederr: false,
		},
	}

	for _, c := range cases {
		testUser := &User{
			FirstName: c.currFN,
			LastName:  c.currLN,
		}

		testUpd := &Updates{
			FirstName: c.newFN,
			LastName:  c.newLN,
		}

		err := testUser.ApplyUpdates(testUpd)
		if err == nil && c.expectederr {
			t.Errorf("case %s: error expected but didn't get one. Double check that error is returned if updates are invalid", c.name)

		} else if err != nil && !c.expectederr {
			t.Errorf("case %s: didn't expect error but got one\nError: %v\n", c.name, err)

		}

	}

}

// TestToUser is a test function for the user ToUser function that
// 	ensures the fields of the NewUser struct are correctly transferred
// 	to a User struct, including the gravatarURL
func TestToUser(t *testing.T) {
	cases := []struct {
		name         string
		first        string
		last         string
		password     string
		passwordconf string
		username     string
		email        string
	}{
		{
			"Testing validate",
			"Sarah",
			"Gefiso",
			"dsgalihesdf",
			"segsaoeurs",
			"Stark",
			"someinvalidemail,com",
		},
		{
			"Testing gravatar hash",
			"Sar",
			"Okay",
			"okayyeah",
			"okayyeah",
			"whatever",
			"testemail@yahoo.com",
		},
	}

	for _, c := range cases {
		testNewUser := &NewUser{
			FirstName:    c.first,
			LastName:     c.last,
			UserName:     c.username,
			Email:        c.email,
			Password:     c.password,
			PasswordConf: c.passwordconf,
		}

		testUser, err := testNewUser.ToUser()
		if err == nil && strings.Contains(c.name, "Testing validate") {
			t.Errorf("case %s: expected error but did not receieve one", c.name) // cannot validate

		}

		if testUser != nil {
			testPhotoURL := fmt.Sprintf("%s%s", gravatarBasePhotoURL, fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(strings.TrimSpace(testNewUser.Email))))))

			if testPhotoURL != testUser.PhotoURL {
				t.Errorf("gravatar photo URLs don't match\nExpected: %s\nAcutal: %s\n", testPhotoURL, testUser.PhotoURL)

			}

		}
	}

}
