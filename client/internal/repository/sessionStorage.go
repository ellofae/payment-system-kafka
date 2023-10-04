package repository

import (
	"github.com/ellofae/payment-system-kafka/config"
	"github.com/gorilla/sessions"
)

var sessionStorage *sessions.CookieStore

func InitSessionStorage(cfg *config.Config) {
	sessionStorage = sessions.NewCookieStore([]byte(cfg.Session.SessionKey))

	sessionStorage.Options.HttpOnly = true
	//sessionStorage.Options.Secure = true
}

func SessionStorage() *sessions.CookieStore {
	return sessionStorage
}
