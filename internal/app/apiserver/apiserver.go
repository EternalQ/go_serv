package apiserver

import (
	"database/sql"
	"go_serv/internal/app/store/sqlstore"
	"net/http"

	"github.com/gorilla/sessions"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	server := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, server)
}

func newDB(dburl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
