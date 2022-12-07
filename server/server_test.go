package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

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
		Context("When inserting a new user", func() {
			var validUserJson string
			var userReader *strings.Reader
			var request *http.Request
			var recorder *httptest.ResponseRecorder

			BeforeEach(func() {
				validUserJson = `{"ID": 12345, "Name": "Josset", "Email": "isset.josset@gmail.com", "Age": 26}`
				userReader = strings.NewReader(validUserJson)
				request = httptest.NewRequest(http.MethodPost, "/user", userReader)
				recorder = httptest.NewRecorder()
			})

			It("Should return a 201 status code if payload is valid", func() {
				server.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(201))
			})

			It("Should return a 405 status code if method is not allowed", func() {
				request := httptest.NewRequest(http.MethodGet, "/user", userReader)
				server.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(405))
			})

			It("Should return a 400 status code when payload is not valid", func() {
				userReader := strings.NewReader(`"{ID": 1234, {}`)
				request := httptest.NewRequest(http.MethodPost, "/user", userReader)
				server.InsertUserHandler(recorder, request)

				result := recorder.Result()
				defer result.Body.Close()
				Expect(result.StatusCode).To(Equal(400))
			})
		})
	})
})
