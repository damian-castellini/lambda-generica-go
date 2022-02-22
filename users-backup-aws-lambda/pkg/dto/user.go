package dto

type User struct {
	Id     string
	Record string
}

func NewUser(record string) User {
	return User{
		Record: record,
	}
}
