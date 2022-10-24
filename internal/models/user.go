package models

import "gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent"

type User struct {
	*ent.User
}

func NewUser(id int64, firstName string, lastName string, userName string) *User {
	return &User{
		User: &ent.User{
			ID:        id,
			FirstName: firstName,
			LastName:  lastName,
			UserName:  userName,
		},
	}
}
