package storage

type User struct {
	id      uint
	name    string
	tickets map[uint]Ticket
}

func newUser(id uint, name string) *User {
	return &User{
		id:   id,
		name: name,
	}
}
