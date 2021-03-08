package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

// NewRouter creates all endpoint for chat app.
func NewRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/postSignup", postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", postSignin).Methods(http.MethodPost)

	router.HandleFunc("/accounts/{id}", getAccount).Methods(http.MethodGet)

	router.HandleFunc("/accounts/{id}/rooms", getAccountRooms).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{id}/rooms", postAccountRooms).Methods(http.MethodPost)
	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}", getAccountRoom).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}", putAccountRoom).Methods(http.MethodPut)

	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}/messages", getMessages).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}/messages", postMessages).Methods(http.MethodPost)

	return router
}

// postSignup handles request for a new account creation.
func postSignup(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// postSignin handles login request for existing user.
func postSignin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// getAccount handles request for user's account information.
func getAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// getAccountRooms returns user's rooms.
func getAccountRooms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// postAccountRooms creates a new room for requesting user.
func postAccountRooms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// getAccountRoom returns room info.
func getAccountRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// putAccountRoom allows to add and remove room members.
func putAccountRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// getMessages returns messages for the selected room.
func getMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// postMessages allows user to create a new message.
func postMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
