package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cl.isset.userfy/model"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if !json.Valid(body) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var id uint = 1
	newUser := model.User{}
	json.Unmarshal(body, &newUser)
	newUser.ID = id
	getUserURL := fmt.Sprintf("/users/%d", id)
	w.Header().Set("Location", getUserURL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}
