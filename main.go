package main

import (
	"fmt"
	"net/http"

	"cl.isset.userfy/database"
	"cl.isset.userfy/repository"
	"cl.isset.userfy/server"
)

func main() {
	db, _ := database.OpenDatabase()
	userRepository := repository.UserRepository{DB: db}
	userServer := server.UserServer{Repository: userRepository}

	http.HandleFunc("/", server.RootHandler)

	http.HandleFunc("/user", userServer.InsertUserHandler)
	http.HandleFunc("/user/update", userServer.UpdateUserHandler)
	http.HandleFunc("/users", userServer.GetUsersHandler)
	http.HandleFunc("/user/delete/", userServer.DeleteUserHandler)

	fmt.Println("Listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
