package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_serv/internal/app/model"
	"go_serv/internal/app/store"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName        = "go-serv"
	ctxUserKey  ctxKey = iota
)

var (
	errIncorrectLoginData = errors.New("incorrect email or password")
	errNotAuthenticated   = errors.New("not authenticated")
)

type ctxKey int16

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionsStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionsStore,
	}

	s.configureRouter()

	return s
}

// implement http.handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "test")
	})

	s.router.HandleFunc("/users", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleCreateSession()).Methods("POST")

	privat := s.router.NewRoute().Subrouter()
	privat.PathPrefix("/private")
	privat.Use(s.authenticateUser)

}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(rw, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(rw, r, http.StatusUnauthorized, errNotAuthenticated)
		}

		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), ctxUserKey, u)))
	})
}

func (s *server) handleCreateUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(rw, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(rw, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.ErasePassword()
		s.respond(rw, r, http.StatusCreated, u)
	}
}

func (s *server) handleCreateSession() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(rw, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(rw, r, http.StatusUnauthorized, errIncorrectLoginData)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, rw, session); err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(rw, r, http.StatusOK, nil)
	}
}

func (s *server) error(rw http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(rw, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(rw http.ResponseWriter, r *http.Request, code int, data interface{}) {
	rw.WriteHeader(code)
	if data != nil {
		json.NewEncoder(rw).Encode(data)
	}
}
