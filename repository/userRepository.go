package repository

import (
	"database/sql"
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
	Exec(query string, args ...any) (sql.Result, error)
}

type UserRepository struct {
	DB SQLDB
}

func (userRepo UserRepository) InsertUser(user model.User) model.User {

	result, err := userRepo.DB.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", user.Name, user.Email, user.Age)
	if err != nil {
		fmt.Printf("error inserting an user: %v\n", err)
		return model.User{}
	}
	lastInsertedId, _ := result.LastInsertId()
	user.ID = uint(lastInsertedId)
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
	userRepo.DB.Exec("UPDATE users SET name = ?, email = ?, age = ?  WHERE ID = ?", user.Name, user.Email, user.Age, user.ID)

	return &user, nil
}

func (userRepo UserRepository) Clear() {
	users = []model.User{}
	idCounter = 0
}
