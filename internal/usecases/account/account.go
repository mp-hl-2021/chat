package account

import (
	"github.com/mp-hl-2021/chat/internal/domain/account"
	"github.com/mp-hl-2021/chat/internal/service/token"

	"golang.org/x/crypto/bcrypt"

	"errors"
	"unicode"
)

var (
	ErrInvalidLoginString    = errors.New("login string contains invalid character")
	ErrInvalidPasswordString = errors.New("password string contains invalid character")
	ErrTooShortString        = errors.New("too short string")
	ErrTooLongString         = errors.New("too long string")

	ErrInvalidLogin    = errors.New("login not found")
	ErrInvalidPassword = errors.New("invalid password")
)

const (
	minLoginLength    = 6
	maxLoginLength    = 32
	minPasswordLength = 14
	maxPasswordLength = 48
)

type Account struct {
	Id string
}

type Interface interface {
	CreateAccount(login, password string) (Account, error)
	GetAccountById(id string) (Account, error)

	LoginToAccount(login, password string) (string, error)
	Authenticate(token string) (string, error)
}

type UseCases struct {
	AccountStorage account.Interface
	Auth           token.Interface
}

func (a *UseCases) CreateAccount(login, password string) (Account, error) {
	if err := validateLogin(login); err != nil {
		return Account{}, err
	}
	if err := validatePassword(password); err != nil {
		return Account{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}
	acc, err := a.AccountStorage.CreateAccount(account.Credentials{
		Login:    login,
		Password: string(hashedPassword),
	})
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, nil
}

func (a *UseCases) GetAccountById(id string) (Account, error) {
	acc, err := a.AccountStorage.GetAccountById(id)
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, err
}

func (a *UseCases) LoginToAccount(login, password string) (string, error) {
	if err := validateLogin(login); err != nil {
		return "", err
	}
	if err := validatePassword(password); err != nil {
		return "", err
	}
	acc, err := a.AccountStorage.GetAccountByLogin(login)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		// todo: check error type. return invalid login or password
		return "", err
	}
	t, err := a.Auth.IssueToken(acc.Id)
	if err != nil {
		return "", err
	}
	return t, err
}

func (a *UseCases) Authenticate(token string) (string, error) {
	return a.Auth.UserIdByToken(token)
}

func validateLogin(login string) error {
	chars := 0
	for _, r := range login {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrInvalidLoginString
		}
		chars++
	}
	if chars < minLoginLength {
		return ErrTooShortString
	}
	if chars > maxLoginLength {
		return ErrTooLongString
	}
	return nil
}

func validatePassword(password string) error {
	chars := 0
	for _, r := range password {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return ErrInvalidPasswordString
		}
		chars++
	}
	if chars < minPasswordLength {
		return ErrTooShortString
	}
	if chars > maxPasswordLength {
		return ErrTooLongString
	}
	return nil
}
