package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/choijos/assignments-choijos/servers/gateway/handlers"
	"github.com/choijos/assignments-choijos/servers/gateway/models/cars"
	"github.com/choijos/assignments-choijos/servers/gateway/models/users"
	"github.com/choijos/assignments-choijos/servers/gateway/sessions"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

// Director is the director used for routing to microservices
type Director func(r *http.Request)

// CustomDirector forwards to the microservice and passes it the current user.
func CustomDirector(targets []*url.URL, ctx *handlers.HandlerContext) Director {
	var counter int32
	counter = 0
	mutex := sync.Mutex{}

	return func(r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		targ := targets[counter%int32(len(targets))]
		atomic.AddInt32(&counter, 1)
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Del("X-User")
		sessionState := &handlers.SessionState{}
		_, err := sessions.GetState(r, ctx.SessKey, ctx.SessStore, sessionState)
		// If there is an error, forward it to the API to deal with it.
		if err != nil {
			r.Header.Add("X-User", "{}")

		} else {
			user := sessionState.AuthUser
			userJSON, err := json.Marshal(user)
			if err != nil {
				r.Header.Add("X-User", "{}")

			} else {
				r.Header.Add("X-User", string(userJSON))

			}

		}

		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme

	}

}

func getURLs(addrString string) []*url.URL {
	addrsSplit := strings.Split(addrString, ",")
	URLs := make([]*url.URL, len(addrsSplit))
	for i, c := range addrsSplit {
		URL, err := url.Parse(c)
		if err != nil {
			log.Fatal(fmt.Printf("Failure to parse url %v", err))

		}

		URLs[i] = URL

	}

	return URLs

}

//main is the main entry point for the server
func main() {
	parkingAddr := os.Getenv("PARKINGADDR")
	if len(parkingAddr) == 0 {
		log.Fatal("No parking address environment variable set")

	}

	sessKey := os.Getenv("SESSIONKEY")
	if len(sessKey) == 0 {
		log.Fatal("No session key environment variable set")

	}

	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		log.Fatal("No redis address environment variable set")

	}

	dsn := os.Getenv("DSN")
	if len(dsn) == 0 {
		log.Fatal("No dsn environment variable set")

	}

	redisDB := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	newrs := sessions.NewRedisStore(redisDB, time.Hour) // sessionDuration?

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Cannot open sql store: %v", err) // Might have to change how we respond
		return

	}

	newCtx := &handlers.HandlerContext{
		SessKey:   sessKey,
		SessStore: newrs,
		UserStore: &users.SQLStore{DbStore: db},
		CarStore:  &cars.SQLStore{DbStore: db},
	}

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"

	}

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	if tlsKeyPath == "" || tlsCertPath == "" { // might have to separate these checks for more tailored error messages
		log.Printf("TLS key or cert not set as environment variable(s): %d \n TLSKEY: %s \n TLSCERT: %s", http.StatusInternalServerError, tlsKeyPath, tlsCertPath)
		return

	}

	// // Create URLs for proxies
	// messagesURLs := getURLs(messagesAddr)
	// summaryURLs := getURLs(summaryAddr)
	// messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(messagesURLs, newCtx)}
	// summaryProxy := &httputil.ReverseProxy{Director: CustomDirector(summaryURLs, newCtx)}
	parkingURLs := getURLs(parkingAddr)
	parkingProxy := &httputil.ReverseProxy{Director: CustomDirector(parkingURLs, newCtx)}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", newCtx.UsersHandler)
	mux.HandleFunc("/v1/users/", newCtx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", newCtx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", newCtx.SpecificSessionHandler)

	// new stuff for assignment
	// r.HandleFunc("/v1/users/{id}/cars", newCtx.UserCarsHandler)
	// r.HandleFunc("/v1/users/{id}/cars/{carid}", newCtx.SpecificUserCarHandler)
	mux.HandleFunc("/v1/cars", newCtx.UserCarsHandler)
	mux.HandleFunc("/v1/cars/", newCtx.SpecificUserCarHandler)

	mux.Handle("/v1/usersparking/", parkingProxy)
	mux.Handle("/v1/parking/", parkingProxy)

	// mux.Handle("/v1/channels", messagesProxy) // double check the round robin stuff
	// mux.Handle("/v1/channels/", messagesProxy)
	// mux.Handle("/v1/messages/", messagesProxy)
	// mux.Handle("/v1/summary", summaryProxy)

	wrappedMux := &handlers.CORS{Handler: mux}

	log.Printf("Server is listening at %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux)) // add tls?

}