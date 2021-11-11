package store

import (
	"go_serv/internal/app/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := r.store.db.
	QueryRow("insert into users (email, encrypted_password) values ($1, $2) returning id", u.Email, u.EncryptedPassword).
	Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	return nil, nil
}
