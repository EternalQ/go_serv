package teststore

import (
	"go_serv/internal/app/model"
	"go_serv/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	v, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return v, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, v := range r.users {
		if v.Email == email {
			return v, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
