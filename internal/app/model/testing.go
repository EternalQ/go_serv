package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "example@example.com",
		Password: "password",
	}
}
