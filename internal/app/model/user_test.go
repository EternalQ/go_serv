package model_test

import (
	"go_serv/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_BeforCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}