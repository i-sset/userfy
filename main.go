package main

import (
	"fmt"
	"net/http"

	"cl.isset.userfy/server"
)

func main() {
	http.HandleFunc("/", server.RootHandler)
	http.HandleFunc("/user", server.InsertUserHandler)
	http.HandleFunc("/users", server.GetUsersHandler)
	fmt.Println("Listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
