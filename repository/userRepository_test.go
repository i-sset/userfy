package repository_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"cl.isset.userfy/model"
	"cl.isset.userfy/repository"
)

var userRepository repository.UserRepository

var _ = Describe("UserRepository", func() {
	Describe("When inserting a valid user", func() {
		var validUser model.User
		var expectedUser model.User

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
		var validUser model.User
		var users []model.User

		BeforeEach(func() {
			validUser = model.User{Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
			userRepository.Clear()
			userRepository.InsertUser(validUser)
			userRepository.InsertUser(validUser)
			users = userRepository.GetUsers()
		})

		It("Should delete all users", func() {
			Expect(len(users)).To(Equal(2))

			userRepository.Clear()
			users = userRepository.GetUsers()
			Expect(len(users)).To(Equal(0))
		})

		It("Should reset the idCounter used for IDs", func() {
			Expect(users[0].ID).To(Equal(uint(1)))
			Expect(users[1].ID).To(Equal(uint(2)))
			userRepository.Clear()

			userRepository.InsertUser(validUser)
			users = userRepository.GetUsers()
			Expect(users[0].ID).To(Equal(uint(1)))

		})
	})

	Describe("When updating an existing user", func() {
		Context("When payload is valid", func() {
			var validUser model.User
			var expectedUser model.User
			BeforeEach(func() {
				validUser = model.User{ID: 1, Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
				expectedUser = model.User{ID: 1, Name: "Joseto", Email: "josset.isset@hotmail.com", Age: 30}
				userRepository.Clear()
			})

			It("Should return the updated entity", func() {
				userRepository.InsertUser(validUser)
				validUser.Name = "Joseto"
				validUser.Email = "josset.isset@hotmail.com"
				validUser.Age = 30

				updatedUser, _ := userRepository.UpdateUser(validUser)

				Expect(updatedUser.ID).To(Equal(expectedUser.ID))
				Expect(updatedUser.Name).To(Equal(expectedUser.Name))
				Expect(updatedUser.Email).To(Equal(expectedUser.Email))
				Expect(updatedUser.Age).To(Equal(expectedUser.Age))
			})
		})

		Context("When payload is not a valid entity", func() {
			var nonExistentUser model.User
			BeforeEach(func() {
				nonExistentUser = model.User{ID: 5, Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
				userRepository.Clear()
			})

			It("Should return an error", func() {
				user, err := userRepository.UpdateUser(nonExistentUser)

				Expect(err).ShouldNot(BeNil())
				Expect(user).Should(BeNil())
			})
		})
	})
})
