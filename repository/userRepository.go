package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"cl.isset.userfy/model"
)

var idCounter uint = 0
var users []model.User

type IUserRepository interface {
	InsertUser(model.User) model.User
	GetUsers() []model.User
	UpdateUser(model.User) (*model.User, error)
}

type SQLDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type UserRepository struct {
	DB SQLDB
}

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

	users := []model.User{}
	usersRow, err := userRepo.DB.Query("SELECT * FROM users")
	if err != nil {
		fmt.Printf("Error while fetching users: %v\n", err)
		return nil
	}

	for usersRow.Next() {
		user := model.User{}
		usersRow.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
		users = append(users, user)
	}
	return users
}

func (userRepo UserRepository) UpdateUser(user model.User) (*model.User, error) {
	var userToBeUpdated *model.User
	for _, u := range users {
		if u.ID == user.ID {
			userToBeUpdated = &u
			break
		}
	}

	if userToBeUpdated == nil {
		return nil, errors.New("user provided does not exist")
	}
	userToBeUpdated.Name = user.Name
	userToBeUpdated.Email = user.Email
	userToBeUpdated.Age = user.Age

	return userToBeUpdated, nil
}

func (userRepo UserRepository) Clear() {
	users = []model.User{}
	idCounter = 0
}
