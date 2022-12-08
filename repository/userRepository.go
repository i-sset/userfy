package repository

import (
	"cl.isset.userfy/model"
)

var idCounter uint = 0
var users []model.User

type UserRepository struct{}

func (userRepo UserRepository) InsertUser(user model.User) model.User {
	id := nextID()
	user.ID = id
	users = append(users, user)
	return user
}

func nextID() uint {
	idCounter = idCounter + uint(1)
	return idCounter
}

func (userRepo UserRepository) GetUsers() []model.User {
	return users
}

func (userRepo UserRepository) Clear() {
	users = []model.User{}
}
