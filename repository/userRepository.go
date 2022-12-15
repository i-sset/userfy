package repository

import (
	"database/sql"
	"fmt"

	"cl.isset.userfy/model"
)

type IUserRepository interface {
	InsertUser(model.User) model.User
	GetUsers() []model.User
	UpdateUser(model.User) (*model.User, error)
	DeleteUser(id int) bool
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

func (userRepo UserRepository) DeleteUser(id int) bool {
	result, err := userRepo.DB.Exec("DELETE FROM users WHERE ID = ?", id)
	if err != nil {
		fmt.Printf("Error while deleting a user: %v\n", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil {
		fmt.Printf("Delete statement Rows affected: %d", rowsAffected)
		return true
	}

	return false
}
