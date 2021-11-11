package store_test

import (
	"go_serv/internal/app/model"
	"go_serv/internal/app/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardow := store.TestStore(t, databaseURL)
	defer teardow("users")

	u, err := s.User().Create(&model.User{
		Email: "example@example.net",
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}