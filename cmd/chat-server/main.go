package main

import (
	"github.com/mp-hl-2021/chat/internal/interface/httpapi"
	"github.com/mp-hl-2021/chat/internal/interface/memory/accountrepo"
	"github.com/mp-hl-2021/chat/internal/interface/memory/messagerepo"
	"github.com/mp-hl-2021/chat/internal/interface/memory/roomrepo"
	"github.com/mp-hl-2021/chat/internal/service/token"
	"github.com/mp-hl-2021/chat/internal/usecases/account"
	"github.com/mp-hl-2021/chat/internal/usecases/message"
	"github.com/mp-hl-2021/chat/internal/usecases/room"

	"flag"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)

	a, err := token.NewJwt(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	accountUseCases := &account.UseCases{
		AccountStorage: accountrepo.NewMemory(),
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
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
