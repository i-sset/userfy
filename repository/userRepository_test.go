package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"cl.isset.userfy/model"
	"cl.isset.userfy/repository"
)

//implements repository.SQLDB interface
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

func (fakeDB *fakeSQLDB) Exec(query string, args ...any) (sql.Result, error) {
	for _, arg := range args {
		argValue := fmt.Sprintf("%v", arg)
		query = strings.Replace(query, "?", argValue, 1)
	}
	fakeDB.Parameter = query

	fakeDB.TimesCalled++
	return fakeResult{}, nil
}

//implements sql.Result interface
type fakeResult struct{}

func (f fakeResult) LastInsertId() (int64, error) {
	return int64(10), nil
}
func (f fakeResult) RowsAffected() (int64, error) {
	return int64(1), nil
}

var fakeDB fakeSQLDB
var userRepository repository.UserRepository

var _ = Describe("UserRepository", func() {
	Describe("When inserting a valid user", func() {
		var validUser model.User

		BeforeEach(func() {
			validUser = model.User{Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
			fakeDB = fakeSQLDB{}
			userRepository = repository.UserRepository{&fakeDB}
		})

		It("Should call db.Exec method once", func() {
			expectedTimes := 1
			userRepository.InsertUser(validUser)

			assertTimesCalled(expectedTimes, fakeDB.TimesCalled)

		})

		It("Should call db.Exec with correct insert query", func() {
			expectedQuery := fmt.Sprintf("INSERT INTO users (name, email, age) VALUES (%s, %s, %d)", validUser.Name, validUser.Email, validUser.Age)
			userRepository.InsertUser(validUser)

			assertQuery(expectedQuery, fakeDB.CalledWith())
		})
	})

	Describe("When fetching all users", func() {
		BeforeEach(func() {
			fakeDB = fakeSQLDB{}
			userRepository = repository.UserRepository{&fakeDB}
		})

		It("Should call db.Query method once", func() {
			expectedTimes := 1
			userRepository.GetUsers()

			assertTimesCalled(expectedTimes, fakeDB.TimesCalled)

		})

		It("Should call db.Query method with 'SELECT * FROM users' query", func() {
			expectedQuery := "SELECT * FROM users"
			userRepository.GetUsers()

			assertQuery(expectedQuery, fakeDB.CalledWith())
		})
	})

	Describe("When updating an existing user", func() {
		Context("When payload is valid", func() {
			var validUser model.User

			BeforeEach(func() {
				validUser = model.User{ID: 1, Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
				fakeDB = fakeSQLDB{}
				userRepository = repository.UserRepository{&fakeDB}
			})

			It("Should call db.Exec method once", func() {
				userRepository.UpdateUser(validUser)
				expectedTimes := 1

				assertTimesCalled(expectedTimes, fakeDB.TimesCalled)
			})

			It("Should call db.Exec method with correct update query", func() {
				expectedQuery := "UPDATE users SET name = Josset, email = isset.josset@gmail.com, age = 26  WHERE ID = 1"
				userRepository.UpdateUser(validUser)

				assertQuery(expectedQuery, fakeDB.CalledWith())
			})
		})
	})

	Describe("When deleting an user", func() {
		Context("When the user exists", func() {
			var validUser model.User

			BeforeEach(func() {
				validUser = model.User{ID: 1, Name: "Josset", Email: "isset.josset@gmail.com", Age: 26}
				fakeDB = fakeSQLDB{}
				userRepository = repository.UserRepository{&fakeDB}
			})

			It("Should call db.Exec method with correct delete query", func() {
				expectedQuery := "DELETE FROM users WHERE ID = 1"
				userRepository.DeleteUser(validUser)

				assertQuery(expectedQuery, fakeDB.CalledWith())
			})

			It("Should call db.Exec method with delete statement once", func() {
				userRepository.DeleteUser(validUser)
				expectedTimes := 1

				assertTimesCalled(expectedTimes, fakeDB.TimesCalled)
			})
		})
	})
})

func assertQuery(expected string, actual string) {
	Expect(actual).To(Equal(expected))
}

func assertTimesCalled(expected int, actual int) {
	Expect(actual).To(Equal(expected))
}
