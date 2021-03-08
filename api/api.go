package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// NewRouter creates all endpoint for chat app.
func NewRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", postSignup).Methods(http.MethodPost)
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

type postSignupRequestModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// postSignup handles request for a new account creation.
func postSignup(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// todo
	location := fmt.Sprintf("/accounts/%s", "1")
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

// postSignin handles login request for existing user.
func postSignin(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// todo
}

type getAccountResponseModel struct {
	Id string `json:"id"`
}

// getAccount handles request for user's account information.
func getAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type getAccountRoomsResponseModel struct {
	RoomIds     []uint64 `json:"room-ids"`
	RoomsNumber uint64   `json:"rooms-number"`
}

// getAccountRooms returns user's rooms.
func getAccountRooms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// postAccountRooms creates a new room for requesting user.
func postAccountRooms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type getAccountRoomResponseModel struct {
	CreatorId    uint64   `json:"creator-id"`
	MemberIds    []uint64 `json:"member-ids"`
	MembersCount uint64   `json:"members-count"`
}

// getAccountRoom returns room info.
func getAccountRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type putAccountRoomRequestModel struct {
	MemberIds []uint64 `json:"member-ids"`
}

// putAccountRoom allows to add and remove room members.
func putAccountRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type getMessagesResponseModel struct {
	Messages []struct {
		AuthorId uint64 `json:"author-id"`
		Text     string `json:"text"`
	} `json:"messages"`
}

// getMessages returns messages for the selected room.
func getMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type postMessagesRequestModel struct {
	Text string `json:"text"`
}

// postMessages allows user to create a new message.
func postMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
