package models

type User struct {
	ID        int64
	FirstName string
	LastName  string
	UserName  string
}

func NewUser(id int64, firstName string, lastName string, userName string) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		UserName:  userName,
	}
}
