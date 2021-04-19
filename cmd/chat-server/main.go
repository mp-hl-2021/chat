package main

import (
	"github.com/mp-hl-2021/chat/internal/interface/httpapi"
	"github.com/mp-hl-2021/chat/internal/interface/memory/messagerepo"
	"github.com/mp-hl-2021/chat/internal/interface/memory/roomrepo"
	"github.com/mp-hl-2021/chat/internal/interface/postgres/accountrepo"
	"github.com/mp-hl-2021/chat/internal/service/token"
	"github.com/mp-hl-2021/chat/internal/usecases/account"
	"github.com/mp-hl-2021/chat/internal/usecases/message"
	"github.com/mp-hl-2021/chat/internal/usecases/room"

	_ "github.com/lib/pq"

	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	if err != nil {
		panic(err)
	}
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)
	if err != nil {
		panic(err)
	}

	a, err := token.NewJwt(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	connStr := "user=postgres password=12345678 host=db dbname=postgres sslmode=disable"

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to DB: %v", err))
	}

	accountUseCases := &account.UseCases{
		AccountStorage: accountrepo.New(conn),
		Auth:           a,
	}
	roomUseCases := &room.UseCases{
		RoomStorage: roomrepo.NewMemory(),
	}
	messageUseCases := &message.UseCases{
		MessageStorage: messagerepo.NewMemory(),
	}

	service := httpapi.NewApi(accountUseCases, roomUseCases, messageUseCases)

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
