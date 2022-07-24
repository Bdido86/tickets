package storage

import "fmt"

var lastUserId = uint(0)

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

func initUser(id uint, name string) *User {
	user := newUser(id, name)
	dataUsers[user.id] = user

	return user
}

func getUserByName(name string) *User {
	for _, user := range dataUsers {
		if name == user.name {
			return user
		}
	}
	return nil
}

func getAuthToken(name string) string {
	user := getUserByName(name)
	if user != nil {
		return getToken(user)
	}

	user = initUser(getNextUserId(), name)
	return getToken(user)
}

func getToken(user *User) string {
	return fmt.Sprintf("%d-%s", user.id, user.name)
}

func getNextUserId() uint {
	for {
		lastUserId++
		if _, ok := dataUsers[lastUserId]; !ok {
			return lastUserId
		}
	}
}
