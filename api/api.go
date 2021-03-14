package api

import (
	"github.com/mp-hl-2021/chat/usecases"

	"github.com/gorilla/mux"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Api struct {
	AccountUseCases usecases.AccountUseCasesInterface
}

func NewApi(a usecases.AccountUseCasesInterface) *Api {
	return &Api{
		AccountUseCases: a,
	}
}

// NewRouter creates all endpoint for chat app.
func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/accounts/{id}", a.authorize(a.getAccount)).Methods(http.MethodGet)

	router.HandleFunc("/accounts/{id}/rooms", a.getAccountRooms).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{id}/rooms", a.postAccountRooms).Methods(http.MethodPost)
	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}", a.getAccountRoom).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}", a.putAccountRoom).Methods(http.MethodPut)

	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}/messages", a.getMessages).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}/rooms/{room_id}/messages", a.postMessages).Methods(http.MethodPost)

	return router
}

type postSignupRequestModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// postSignup handles request for a new account creation.
func (a *Api) postSignup(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := a.AccountUseCases.CreateAccount(m.Login, m.Password)
	if err != nil { // todo: map domain errors to http error codes
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	location := fmt.Sprintf("/accounts/%s", acc.Id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

// postSignin handles login request for existing user.
func (a *Api) postSignin(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.AccountUseCases.LoginToAccount(m.Login, m.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	w.Write([]byte(token))
}

type getAccountResponseModel struct {
	Id string `json:"id"`
}

// getAccount handles request for user's account information.
func (a *Api) getAccount(w http.ResponseWriter, r *http.Request) {
	accountId, ok := r.Context().Value("account_id").(string) // todo: move to constant
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id != accountId { // todo: need an authorization module
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	acc, err := a.AccountUseCases.GetAccountById(accountId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := getAccountResponseModel{Id: acc.Id}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type getAccountRoomsResponseModel struct {
	RoomIds     []uint64 `json:"room-ids"`
	RoomsNumber uint64   `json:"rooms-number"`
}

// getAccountRooms returns user's rooms.
func (a *Api) getAccountRooms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// postAccountRooms creates a new room for requesting user.
func (a *Api) postAccountRooms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type getAccountRoomResponseModel struct {
	CreatorId    string   `json:"creator-id"`
	MemberIds    []string `json:"member-ids"`
	MembersCount uint64   `json:"members-count"`
}

// getAccountRoom returns room info.
func (a *Api) getAccountRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type putAccountRoomRequestModel struct {
	MemberIds []string `json:"member-ids"`
}

// putAccountRoom allows to add and remove room members.
func (a *Api) putAccountRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type getMessagesResponseModel struct {
	Messages []struct {
		AuthorId string `json:"author-id"`
		Text     string `json:"text"`
	} `json:"messages"`
}

// getMessages returns messages for the selected room.
func (a *Api) getMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

type postMessagesRequestModel struct {
	Text string `json:"text"`
}

// postMessages allows user to create a new message.
func (a *Api) postMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *Api) authorize(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearHeader := r.Header.Get("Authorization")
		strArr := strings.Split(bearHeader, " ")
		if len(strArr) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := strArr[1]
		id, err := a.AccountUseCases.Authenticate(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "account_id", id)
		handler(w, r.WithContext(ctx))
	}
}
