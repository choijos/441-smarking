package handlers

import (
	"time"

	"github.com/choijos/assignments-choijos/servers/gateway/models/users"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

type SessionState struct {
	StartTime time.Time  `json:"time"`
	AuthUser  *users.User `json:"authUser"` // maybe pointer
}
