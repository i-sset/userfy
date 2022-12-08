package main

import (
	"net/http"

	"cl.isset.userfy/server"
)

func main() {
	http.HandleFunc("/", server.RootHandler)
	http.HandleFunc("/user", server.InsertUserHandler)
	http.ListenAndServe(":3000", nil)
}
