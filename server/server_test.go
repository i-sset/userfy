package server_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"cl.isset.userfy/model"
	"cl.isset.userfy/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type fakeUserRepository struct {
	timesCalled int
	idParameter int
}

func (userRepo *fakeUserRepository) InsertUser(user model.User) model.User {
	return user
}

func (userRepo *fakeUserRepository) GetUsers() []model.User {
	return []model.User{
		model.User{},
	}
}

func (userRepo *fakeUserRepository) UpdateUser(user model.User) (*model.User, error) {
	notExistentID := uint(14)

	if user.ID == notExistentID {
		return nil, errors.New("user provided does not exist")
	}
	return &user, nil
}

func (userRepo *fakeUserRepository) DeleteUser(id int) bool {
	notExistentID := 14
	userRepo.timesCalled++
	userRepo.idParameter = id

	return id != notExistentID
}

var fakeRepo fakeUserRepository
var userServer server.UserServer

var _ = Describe("Server", func() {
	Describe("root endpoint", func() {
		Context("When everything is ok", func() {
			It("Should return 200 status code", func() {
				request := httptest.NewRequest(http.MethodGet, "/", nil)
				writer := httptest.NewRecorder()
				server.RootHandler(writer, request)

				result := writer.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(200))
			})
		})
	})

	Describe("/user endpoint", func() {
		var validUserJson string
		var userReader *strings.Reader
		var request *http.Request
		var recorder *httptest.ResponseRecorder

		Context("When inserting a new valid user", func() {
			BeforeEach(func() {
				validUserJson = `{"Name": "Josset", "Email": "isset.josset@gmail.com", "Age": 26}`
				userReader = strings.NewReader(validUserJson)
				request = httptest.NewRequest(http.MethodPost, "/user", userReader)
				recorder = httptest.NewRecorder()
				fakeRepo = fakeUserRepository{}
				userServer = server.UserServer{&fakeRepo}
			})

			It("Should return a 201 status code if payload is valid", func() {
				userServer.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(http.StatusCreated))
			})

			It("Should provide a link to the new user created in the Location header", func() {
				userServer.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				createdUser := model.User{}
				json.NewDecoder(result.Body).Decode(&createdUser)

				createdUserURL := fmt.Sprintf("/users/%d", createdUser.ID)
				location := result.Header.Get("Location")
				Expect(location).To(Equal(createdUserURL))
			})
		})

		Context("When inserting a no valid user", func() {
			BeforeEach(func() {
				recorder = httptest.NewRecorder()
			})

			It("Should return a 405 status code if method is not allowed", func() {
				request := httptest.NewRequest(http.MethodGet, "/user", nil)
				userServer.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})

			It("Should return a 400 status code when payload is not valid", func() {
				userReader := strings.NewReader(`"{ID": 1234, {}`)
				request := httptest.NewRequest(http.MethodPost, "/user", userReader)
				userServer.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("/user/update endpoint: update users", func() {
		var endpoint = "/user/update"
		var userReader *strings.Reader
		var request *http.Request
		var recorder *httptest.ResponseRecorder
		Context("When it is a valid request", func() {
			var validUserJson string

			BeforeEach(func() {
				validUserJson = `{"ID": 1,"Name": "Joseto", "Email": "josset.isset@hotmail.com", "Age": 30}`
				userReader = strings.NewReader(validUserJson)
				request = httptest.NewRequest(http.MethodPut, endpoint, userReader)
				recorder = httptest.NewRecorder()
				fakeRepo = fakeUserRepository{}
				userServer = server.UserServer{&fakeRepo}
			})

			It("Should return 200 status code", func() {
				userServer.UpdateUserHandler(recorder, request)

				result := recorder.Result()
				Expect(result.StatusCode).To(Equal(http.StatusOK))
			})

			It("Should return an user entity ", func() {
				userServer.UpdateUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				user := model.User{}
				json.NewDecoder(result.Body).Decode(&user)
				Expect(user).ToNot(HaveField("ID", uint(0)))
			})
		})

		Context("When it is not a valid request", func() {
			var notValidUser string
			var validButNotExistentUser string
			BeforeEach(func() {
				notValidUser = `{"ID": 1,"Nae": "Joseto", "Email: "josset.isset@hotmail.com", "Age" 30}`
				validButNotExistentUser = `{"ID": 14,"Name": "Josetqqweip", "Email": "josset.isset@hotmail.com", "Age": 25}`
				recorder = httptest.NewRecorder()
				fakeRepo = fakeUserRepository{}
				userServer = server.UserServer{&fakeRepo}
			})

			It("Should return a bad request response for a not valid json", func() {
				userReader = strings.NewReader(notValidUser)
				request = httptest.NewRequest(http.MethodPut, endpoint, userReader)
				userServer.UpdateUserHandler(recorder, request)

				result := recorder.Result()

				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("Should return a not found response for a not existent user", func() {
				userReader = strings.NewReader(validButNotExistentUser)
				request = httptest.NewRequest(http.MethodPut, endpoint, userReader)
				userServer.UpdateUserHandler(recorder, request)

				result := recorder.Result()

				Expect(result.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("/users endpoint: get users", func() {
		var request *http.Request
		var recorder *httptest.ResponseRecorder

		BeforeEach(func() {
			request = httptest.NewRequest(http.MethodGet, "/users", nil)
			recorder = httptest.NewRecorder()
			fakeRepo = fakeUserRepository{}
			userServer = server.UserServer{&fakeRepo}
		})

		It("Should return a 200 status code", func() {
			userServer.GetUsersHandler(recorder, request)

			result := recorder.Result()
			Expect(result.StatusCode).To(Equal(http.StatusOK))
		})

		It("Should return a not empty array", func() {
			userServer.GetUsersHandler(recorder, request)

			result := recorder.Result()
			var users []model.User
			json.NewDecoder(result.Body).Decode(&users)
			Expect(users).ShouldNot(BeEmpty())
		})

		Context("When Method GET is not used for this endpoint", func() {
			It("Should return a method not allowed status", func() {
				request = httptest.NewRequest(http.MethodPut, "/users", nil)
				recorder = httptest.NewRecorder()
				userServer.GetUsersHandler(recorder, request)

				result := recorder.Result()
				Expect(result.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})

	Describe("/users/delete endpoint: delete user", func() {
		Context("When the user exists", func() {
			var request *http.Request
			var recorder *httptest.ResponseRecorder
			BeforeEach(func() {
				request = httptest.NewRequest(http.MethodDelete, "/user/delete/1", nil)
				recorder = httptest.NewRecorder()
				fakeRepo = fakeUserRepository{}
				userServer = server.UserServer{&fakeRepo}
			})

			It("Should return a 204 status code", func() {
				userServer.DeleteUserHandler(recorder, request)

				result := recorder.Result()
				Expect(result.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("Should call repository.DeleteUser method once", func() {
				userServer.DeleteUserHandler(recorder, request)

				Expect(fakeRepo.timesCalled).To(Equal(1))
			})

			It("Should call repository.DeleteUser with correct parameter", func() {
				existentId := 1
				userServer.DeleteUserHandler(recorder, request)

				Expect(fakeRepo.idParameter).To(Equal(existentId))
			})
		})

		Context("When the user  does not exists", func() {
			var request *http.Request
			var recorder *httptest.ResponseRecorder
			BeforeEach(func() {
				request = httptest.NewRequest(http.MethodDelete, "/user/delete/14", nil)
				recorder = httptest.NewRecorder()
				fakeRepo = fakeUserRepository{}
				userServer = server.UserServer{&fakeRepo}
			})

			It("Should return a not found status", func() {
				userServer.DeleteUserHandler(recorder, request)

				result := recorder.Result()
				Expect(result.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("When the id provided is not valid", func() {
			var request *http.Request
			var recorder *httptest.ResponseRecorder

			BeforeEach(func() {
				request = httptest.NewRequest(http.MethodDelete, "/user/delete/notValidID", nil)
				recorder = httptest.NewRecorder()
				fakeRepo = fakeUserRepository{}
				userServer = server.UserServer{&fakeRepo}
			})
			It("Should return a bad request response", func() {
				userServer.DeleteUserHandler(recorder, request)

				result := recorder.Result()
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the handler is invoked", func() {
			var request *http.Request
			var recorder *httptest.ResponseRecorder

			BeforeEach(func() {
				request = httptest.NewRequest(http.MethodGet, "/user/delete/notValidID", nil)
				recorder = httptest.NewRecorder()
				fakeRepo = fakeUserRepository{}
				userServer = server.UserServer{&fakeRepo}
			})
			It("Should return not allowed response if http method is different of DELETE", func() {
				userServer.DeleteUserHandler(recorder, request)

				result := recorder.Result()
				Expect(result.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})

})
