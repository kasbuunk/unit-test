package repo

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"

	"github.com/kasbuunk/unit-test/repository/storage"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx   context.Context
	repo  Repository
	users []*User
}

func (s *RepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	dbConfig := storage.Config{
		Host:    "localhost",
		Port:    5432,
		User:    "postgres",
		Pass:    "postgres",
		Name:    "unit_test",
		SSLMode: "disable",
	}
	db, err := storage.Connect(dbConfig)
	s.Require().NoError(err, "connecting to database")

	s.repo = NewRepository(db)

	goose.SetTableName("goose_auth")
	err = goose.Run("up", db, "migration")
	s.Require().NoError(err, "running migrations")

	_, err = s.repo.UserDeleteAll(s.ctx)
	s.NoError(err)
}

func (s *RepositoryTestSuite) SetupTest() {
	usrs := []*User{
		{
			Email:        EmailAddress("user@example.com"),
			PasswordHash: PasswordHash("klasjdflksjdklfj"),
		},
		{
			Email:        EmailAddress("another-user@example.com"),
			PasswordHash: PasswordHash("rinsTIErnsietrmies"),
		},
	}

	for _, usr := range usrs {
		_, err := s.repo.UserSave(s.ctx, *usr)
		s.Require().NoError(err)
	}
	savedUsers, err := s.repo.Users(s.ctx)
	s.Require().NoError(err)

	s.users = []*User{}
	s.users = append(s.users, savedUsers...)
}

func (s *RepositoryTestSuite) TearDownTest() {
	_, err := s.repo.UserDeleteAll(s.ctx)
	s.NoError(err)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestUsersCreate() {
	usr := User{
		Email:        EmailAddress("yetanotheruser@example.com"),
		PasswordHash: PasswordHash("456|4{}56rmsteirsnteir"),
	}
	savedUser, err := s.repo.UserSave(s.ctx, usr)
	s.Require().NoError(err)
	s.NotEqual("", savedUser.ID.String())

	rowsAffected, err := s.repo.UserDelete(s.ctx, savedUser.ID)
	s.Require().NoError(err)
	s.NotEqual(rowsAffected, usr)
}

func (s *RepositoryTestSuite) TestUsersSave() {
	existingUsers, err := s.repo.Users(s.ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(existingUsers)
	existingUser := existingUsers[0]
	s.Require().NotEmpty(existingUser.ID)

	fetchedUser, err := s.repo.User(s.ctx, existingUser.ID)
	s.Require().NoError(err)
	s.Equal(fetchedUser, existingUser)

	newEmail := EmailAddress("new@example.com")
	existingUser.Email = EmailAddress(newEmail)
	changedUser, err := s.repo.UserSave(s.ctx, *existingUser)
	s.Require().NoError(err)
	s.Equal(changedUser.Email, newEmail)
}

func (s *RepositoryTestSuite) TestUsersDelete() {
	usr := s.users[0]
	s.Require().NotEmpty(usr.ID)
	fetchedUser, err := s.repo.User(s.ctx, usr.ID)
	s.Require().NoError(err)
	rowsAffected, err := s.repo.UserDelete(s.ctx, fetchedUser.ID)
	s.Require().NoError(err)
	s.Equal(1, rowsAffected)
	deletedUser, err := s.repo.UserDelete(s.ctx, usr.ID)
	s.Require().NoError(err)
	s.NotEqual(deletedUser, usr)
}
