package repository_test

import (
	"database/sql"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"cl.isset.userfy/model"
	"cl.isset.userfy/repository"
)

type fakeSQLDB struct {
	TimesCalled int
	Parameter   string
}

func (fakeDB *fakeSQLDB) CalledWith() string {
	return fakeDB.Parameter
}

func (fakeDB *fakeSQLDB) Query(parameter string, v ...any) (*sql.Rows, error) {
	fakeDB.Parameter = parameter
	fakeDB.TimesCalled++
	return nil, errors.New("Wrong query")
}

var fakeDB = fakeSQLDB{}
var userRepository = repository.UserRepository{&fakeDB}

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
		It("Should call db.Query method once", func() {
			userRepository.GetUsers()
			timesCalled := fakeDB.TimesCalled
			Expect(timesCalled).To(Equal(1))
		})

		It("Should call db.Query method with 'SELECT * FROM users' query", func() {
			userRepository.GetUsers()
			queryCalledWith := fakeDB.CalledWith()

			Expect(queryCalledWith).To(Equal("SELECT * FROM users"))
		})
	})

	Describe("When updating an existing user", func() {
		Context("When payload is valid", func() {
			var validUser model.User
			var expectedUser model.User
			BeforeEach(func() {
				validUser = model.User{ID: 1, Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
				expectedUser = model.User{ID: 1, Name: "Joseto", Email: "josset.isset@hotmail.com", Age: 30}
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
			})

			It("Should return an error", func() {
				user, err := userRepository.UpdateUser(nonExistentUser)

				Expect(err).ShouldNot(BeNil())
				Expect(user).Should(BeNil())
			})
		})
	})
})
