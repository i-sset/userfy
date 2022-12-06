package server_test

import (
	"net/http"
	"net/http/httptest"

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

})
