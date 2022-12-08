package repository_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"cl.isset.userfy/model"
	"cl.isset.userfy/repository"
)

var _ = Describe("UserRepository", func() {
	Describe("When inserting a valid user", func() {
		var validUser model.User
		var expectedUser model.User
		var userRepository repository.UserRepository

		BeforeEach(func() {
			validUser = model.User{Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
			expectedUser = model.User{ID: 1, Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
		})

		It("Should return the user with a new unique id", func() {
			createdUser := userRepository.InsertUser(validUser)

			Expect(createdUser.ID).To(Equal(expectedUser.ID))

			createdUser = userRepository.InsertUser(validUser)
			expectedUser = model.User{ID: 2, Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}

			Expect(createdUser.ID).To(Equal(expectedUser.ID))
		})

		It("Should return the user with the same fields", func() {
			createdUser := userRepository.InsertUser(validUser)

			Expect(createdUser.Name).To(Equal(expectedUser.Name))
			Expect(createdUser.Email).To(Equal(expectedUser.Email))
			Expect(createdUser.Age).To(Equal(expectedUser.Age))
		})
	})

	Describe("When fetching all users", func() {
		var userRepository repository.UserRepository
		var validUser model.User

		BeforeEach(func() {
			validUser = model.User{Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
			userRepository.Clear()
		})

		It("Should return a slice with a length of 2", func() {
			userRepository.InsertUser(validUser)
			userRepository.InsertUser(validUser)
			users := userRepository.GetUsers()

			Expect(len(users)).To(Equal(2))
		})
	})

	Describe("When clear method is called", func() {
		var userRepository repository.UserRepository
		var validUser model.User

		BeforeEach(func() {
			validUser = model.User{Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
			userRepository.Clear()
		})

		It("Should delete all users", func() {
			userRepository.InsertUser(validUser)
			userRepository.InsertUser(validUser)
			users := userRepository.GetUsers()

			Expect(len(users)).To(Equal(2))

			userRepository.Clear()
			users = userRepository.GetUsers()
			Expect(len(users)).To(Equal(0))
		})
	})
})
