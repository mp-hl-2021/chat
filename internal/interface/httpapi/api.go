package httpapi

import (
	"github.com/mp-hl-2021/chat/internal/usecases/account"
	"github.com/mp-hl-2021/chat/internal/usecases/message"
	"github.com/mp-hl-2021/chat/internal/usecases/room"

	"github.com/gorilla/mux"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	accountIdContextKey = "account_id"
	accountIdUrlPathKey = "account_id"
	roomsIdUrlPathKey   = "room_id"
)

type Api struct {
	AccountUseCases account.Interface
	RoomUseCases    room.Interface
	MessageUseCases message.Interface
}

func NewApi(a account.Interface, r room.Interface, m message.Interface) *Api {
	return &Api{
		AccountUseCases: a,
		RoomUseCases:    r,
		MessageUseCases: m,
	}
}

// NewRouter creates all endpoint for chat app.
func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/accounts/"+accountIdUrlPathKey, a.authenticate(a.getAccount)).Methods(http.MethodGet)

	router.HandleFunc("/rooms", a.authenticate(a.getAccountRooms)).Methods(http.MethodGet)
	router.HandleFunc("/rooms", a.authenticate(a.postAccountRooms)).Methods(http.MethodPost)
	router.HandleFunc("/rooms/"+roomsIdUrlPathKey, a.authenticate(a.getAccountRoom)).Methods(http.MethodGet)
	router.HandleFunc("/rooms/"+roomsIdUrlPathKey, a.authenticate(a.putAccountRoom)).Methods(http.MethodPut)

	router.HandleFunc("/rooms/"+roomsIdUrlPathKey+"/messages", a.authenticate(a.getMessages)).Methods(http.MethodGet)
	router.HandleFunc("/rooms/"+roomsIdUrlPathKey+"/messages", a.authenticate(a.postMessages)).Methods(http.MethodPost)

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
	accountId, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	id, ok := vars[accountIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id != accountId { // todo: move to the upper layer
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
	RoomIds     []string `json:"room-ids"`
	RoomsNumber int      `json:"rooms-number"`
}

// getAccountRooms returns user's rooms.
func (a *Api) getAccountRooms(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rr, err := a.RoomUseCases.ListRooms(aid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := getAccountRoomsResponseModel{
		RoomIds:     make([]string, 0, len(rr)),
		RoomsNumber: len(rr),
	}
	for i := range rr {
		m.RoomIds[i] = rr[i].Id
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// postAccountRooms creates a new room for requesting user.
func (a *Api) postAccountRooms(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rm, err := a.RoomUseCases.CreateRoom(aid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	location := fmt.Sprintf("/rooms/%s", rm.Id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

type getAccountRoomResponseModel struct {
	CreatorId    string   `json:"creator-id"`
	MemberIds    []string `json:"member-ids"`
	MembersCount int      `json:"members-count"`
}

// getAccountRoom returns room info.
func (a *Api) getAccountRoom(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	rid, ok := vars[roomsIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rm, err := a.RoomUseCases.GetRoomById(aid, rid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := getAccountRoomResponseModel{
		CreatorId:    "todo", // todo
		MemberIds:    rm.Members,
		MembersCount: len(rm.Members),
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type putAccountRoomRequestModel struct {
	Members []struct {
		Id     string `json:"id"`
		Delete bool   `json:"delete"`
	} `json:"members"`
}

// putAccountRoom allows to add and remove room members.
func (a *Api) putAccountRoom(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	rid, ok := vars[roomsIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := putAccountRoomRequestModel{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// note: this action is made by two transaction within one request.
	// it may lead to a strange situation where add accounts transaction succeeded,
	// but delete accounts transaction failed and whole put request failed with state change.
	// so, here is a good place to change API or use cases and repository.
	toDelete := make([]string, 0, len(m.Members))
	toAdd := make([]string, 0, len(m.Members))
	for _, member := range m.Members {
		if member.Delete {
			toDelete = append(toDelete, member.Id)
			continue
		}
		toAdd = append(toAdd, member.Id)
	}
	err := a.RoomUseCases.AddMembers(aid, rid, toAdd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = a.RoomUseCases.RemoveMembers(aid, rid, toDelete)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type getMessagesResponseModel struct {
	Messages []messageModel `json:"messages"`
}

type messageModel struct {
	AuthorId string `json:"author-id"`
	Text     string `json:"text"`
}

// getMessages returns messages for the selected room.
func (a *Api) getMessages(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	rid, ok := vars[roomsIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgs, err := a.MessageUseCases.ListMessages(aid, rid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := getMessagesResponseModel{Messages: make([]messageModel, 0, len(msgs))}
	for _, msg := range msgs {
		m.Messages = append(m.Messages, messageModel{
			AuthorId: msg.Author,
			Text:     msg.Text,
		})
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type postMessagesRequestModel struct {
	Text string `json:"text"`
}

// postMessages allows user to create a new message.
func (a *Api) postMessages(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value(accountIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	rid, ok := vars[roomsIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var m postMessagesRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.MessageUseCases.CreateMessage(aid, rid, m.Text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) authenticate(handler http.HandlerFunc) http.HandlerFunc {
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
		ctx := context.WithValue(r.Context(), accountIdContextKey, id)
		handler(w, r.WithContext(ctx))
	}
}
