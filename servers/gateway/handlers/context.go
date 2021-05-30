package handlers

import (
	"github.com/choijos/assignments-choijos/servers/gateway/models/cars"
	"github.com/choijos/assignments-choijos/servers/gateway/models/users"
	"github.com/choijos/assignments-choijos/servers/gateway/sessions"
)

// Handler context that is a receiver on any of your HTTP
//	handler functions that need access to globals
type HandlerContext struct {
	SessKey   string
	SessStore sessions.Store
	UserStore users.Store
	CarStore  cars.Store
}
