package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cl.isset.userfy/repository"

	"cl.isset.userfy/model"
)

var userRepository = repository.UserRepository{}

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

	newUser := model.User{}
	json.Unmarshal(body, &newUser)
	createdUser := userRepository.InsertUser(newUser)

	createdUserURL := fmt.Sprintf("/users/%d", createdUser.ID)
	w.Header().Set("Location", createdUserURL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createdUser)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	users := userRepository.GetUsers()
	json.NewEncoder(w).Encode(users)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if !json.Valid(body) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	incomingUser := model.User{}
	json.Unmarshal(body, &incomingUser)
	updatedUser, err := userRepository.UpdateUser(incomingUser)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}
