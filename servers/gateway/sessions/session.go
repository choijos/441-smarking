package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	// "encoding/json"
	"errors"
	"fmt"
	"net/http"
	// "github.com/patrickmn/go-cache"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")


//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//TODO:
	//- create a new SessionID
	//- save the sessionState to the store
	//- add a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	//  (note the constants declared for you above, which will help you avoid typos)
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("signing key may not be empty")

	}

	randomBytes := make([]byte, idLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return InvalidSessionID, err

	}
	remaining := hmac.New(sha256.New, []byte(signingKey))
	remaining.Write(randomBytes)
	remainingBytes := remaining.Sum(nil)
	// https://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go
	finalByteSlice := append(randomBytes, remainingBytes...)
	// Encode the byteslice to Base64 URL Encoded string
	finalSessionID := SessionID(base64.URLEncoding.EncodeToString(finalByteSlice))

	store.Save(finalSessionID, sessionState)

	// authHeader := headerAuthorization + ": " + schemeBearer + " "

	// w.Header().Add(authHeader, string(finalSessionID))
	w.Header().Add(headerAuthorization, fmt.Sprintf("%s%s", schemeBearer, string(finalSessionID)))

	return finalSessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.
	// authHeaderVal := r.Header.Get((headerAuthorization + ": " + schemeBearer + " "))
	authHeaderVal := r.Header.Get(headerAuthorization) // Bearer <SessionID>
	if len(authHeaderVal) == 0 {
		val, valok := r.URL.Query()["auth"]
		if !valok || len(val[0]) < 1 {
			return InvalidSessionID, ErrInvalidScheme

		}

		authHeaderVal = val[0]

	}

	if authHeaderVal == "" {
		return InvalidSessionID, ErrNoSessionID

	}

	authParts := strings.Split(authHeaderVal, " ")

	if authParts[0] != "Bearer" {
		return InvalidSessionID, ErrInvalidScheme

	}

	decodedID, err := base64.URLEncoding.DecodeString(authParts[1])
	if err != nil {
		return InvalidSessionID, err

	}

	idPortion := decodedID[0:idLength]
	compare := decodedID[idLength:]
	remaining := hmac.New(sha256.New, []byte(signingKey))
	_, writeErr := remaining.Write(idPortion)
	if writeErr != nil {
		return InvalidSessionID, writeErr

	}

	remainingBytes := remaining.Sum(nil)
	if hmac.Equal(compare, remainingBytes) {
		return SessionID(authParts[1]), nil

	}

	// return InvalidSessionID, ErrInvalidID

	return InvalidSessionID, fmt.Errorf("Signing key / session id is not valid\n")

}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//TODO: get the SessionID from the request, and get the data
	//associated with that SessionID from the store.
	sessID := r.Header.Get(headerAuthorization)

	if len(sessID) == 0 {
		return SessionID(sessID), ErrNoSessionID

	}

	idParts := strings.Split(sessID, " ")

	if cap(idParts) != 1 {
		sessID = idParts[1]

	} else {
		sessID = idParts[0]

	}

	err := store.Get(SessionID(sessID), sessionState)

	if err != nil {
		return SessionID(sessID), ErrStateNotFound

	}

	return SessionID(sessID), nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//TODO: get the SessionID from the request, and delete the
	//data associated with it in the store.
	sessID := r.Header.Get(headerAuthorization)
	// fmt.Println(sessID)

	if len(sessID) == 0 {
		return SessionID(sessID), ErrNoSessionID

	}

		idParts := strings.Split(sessID, " ")

	if cap(idParts) != 1 {
		sessID = idParts[1]

	} else {
		sessID = idParts[0]

	}

	store.Delete(SessionID(sessID))

	return SessionID(sessID), nil
}
