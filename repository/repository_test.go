package repo

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"

	"github.com/kasbuunk/unit-test/repository/storage"
)

type RepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo Repository
}

func (s *RepositoryTestSuite) SetupSuite() {
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

	s.db = db
	s.repo = NewRepository(s.db)

	dir, err := os.Getwd()
	s.Require().NoError(err, "getting working directory")

	goose.SetTableName("goose_auth")
	err = goose.Run("up", db, filepath.Join(dir, "migration"))
	s.Require().NoError(err, "running migrations")

	_, err = s.repo.UserDeleteAll(context.Background())
	s.NoError(err)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestUsersCRUD() {
	ctx := context.Background()
	users, err := s.repo.Users(ctx)
	s.Require().NoError(err)
	s.Empty(users)

	usr := User{
		Email:        EmailAddress("user@example.com"),
		PasswordHash: PasswordHash("klasjdflksjdklfj"),
	}
	savedUser, err := s.repo.UserSave(ctx, usr)
	s.Require().NoError(err)
	s.NotEqual("", savedUser.ID.String())

	fetchedUser, err := s.repo.User(ctx, savedUser.ID)
	s.Require().NoError(err)
	s.Equal(fetchedUser, savedUser)

	newEmail := EmailAddress("new@example.com")
	savedUser.Email = EmailAddress(newEmail)
	changedUser, err := s.repo.UserSave(ctx, *savedUser)
	s.Require().NoError(err)
	s.Equal(changedUser.Email, newEmail)
}
