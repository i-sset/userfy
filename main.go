package main

import (
	"fmt"
	"net/http"

	"cl.isset.userfy/model"
	"cl.isset.userfy/repository"

	"cl.isset.userfy/server"
)

func main() {
	http.HandleFunc("/", server.RootHandler)
	http.HandleFunc("/user", server.InsertUserHandler)
	http.HandleFunc("/user/update", server.UpdateUserHandler)
	http.HandleFunc("/users", server.GetUsersHandler)

	loadDatabase()
	fmt.Println("Listening at port 8080")
	http.ListenAndServe(":8080", nil)
}

func loadDatabase() {
	userRepository := repository.UserRepository{}

	userRepository.InsertUser(model.User{Name: "Josset Garcia", Email: "isset.joset@gmail.com", Age: 26})
	userRepository.InsertUser(model.User{Name: "Silvana Ferreiro", Email: "silvanaf@thoughtworks.com", Age: 42})
	userRepository.InsertUser(model.User{Name: "Nicolas Bedregal", Email: "nicobe@gmail.com", Age: 32})
	userRepository.InsertUser(model.User{Name: "Javiera Lasus", Email: "javivu@thoughtworks.com", Age: 27})
	userRepository.InsertUser(model.User{Name: "Adriana Ortega", Email: "adriortega@gmail.com", Age: 30})

}
