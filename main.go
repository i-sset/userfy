package main

import (
	"fmt"
	"net/http"

	"cl.isset.userfy/server"
)

func main() {
	http.HandleFunc("/", server.RootHandler)
	http.HandleFunc("/user", server.InsertUserHandler)
	fmt.Println("Listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
