package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"cl.isset.userfy/model"
	"cl.isset.userfy/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			})

			It("Should return a 201 status code if payload is valid", func() {
				server.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(http.StatusCreated))
			})

			It("Should provide a link to the new user created in the Location header", func() {
				server.InsertUserHandler(recorder, request)

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
				server.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})

			It("Should return a 400 status code when payload is not valid", func() {
				userReader := strings.NewReader(`"{ID": 1234, {}`)
				request := httptest.NewRequest(http.MethodPost, "/user", userReader)
				server.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("/users endpoint: get users", func() {
		var request *http.Request
		var recorder *httptest.ResponseRecorder

		BeforeEach(func() {
			request = httptest.NewRequest(http.MethodGet, "/users", nil)
			recorder = httptest.NewRecorder()
		})

		It("Should return a 200 status code", func() {
			server.GetUsersHandler(recorder, request)

			result := recorder.Result()
			Expect(result.StatusCode).To(Equal(http.StatusOK))
		})

		It("Should return a not empty array", func() {
			server.GetUsersHandler(recorder, request)

			result := recorder.Result()
			var users []model.User
			json.NewDecoder(result.Body).Decode(&users)
			Expect(users).ShouldNot(BeEmpty())
		})

		Context("When Method GET is not used for this endpoint", func() {
			It("Should return a method not allowed status", func() {
				request = httptest.NewRequest(http.MethodPut, "/users", nil)
				recorder = httptest.NewRecorder()
				server.GetUsersHandler(recorder, request)

				result := recorder.Result()
				Expect(result.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})
})
