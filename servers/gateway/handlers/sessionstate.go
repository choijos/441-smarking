package handlers

import (
	"time"

	"github.com/choijos/assignments-choijos/servers/gateway/models/users"
)

// Session state struct
type SessionState struct {
	StartTime time.Time  `json:"time"`
	AuthUser  *users.User `json:"authUser"`
}
