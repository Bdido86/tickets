package storage

type User struct {
	id      uint
	name    string
	tickets map[uint]*Ticket
}

func newUser(id uint, name string) *User {
	tickets := make(map[uint]*Ticket)
	return &User{
		id:      id,
		name:    name,
		tickets: tickets,
	}
}
