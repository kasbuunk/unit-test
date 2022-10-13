package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/kasbuunk/unit-test/repository/gen/pgmodel"
)

type User struct {
	ID           uuid.UUID
	Email        EmailAddress
	PasswordHash PasswordHash
}

type (
	EmailAddress string
	PasswordHash string
	Password     string
)

// Repository implements the UserRepository interface for the User
type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (r Repository) User(ctx context.Context, id uuid.UUID) (*User, error) {
	pgUser, err := pgmodel.Users(pgmodel.UserWhere.ID.EQ(id.String())).One(ctx, r.DB)
	if err != nil {
		return nil, errors.Wrap(err, "getting user")
	}
	usr := userToDomainModel(pgUser)
	return usr, nil
}

func (r Repository) Users(ctx context.Context) ([]*User, error) {
	pgUsers, err := pgmodel.Users().All(ctx, r.DB)
	if err != nil {
		return nil, errors.Wrap(err, "getting all users")
	}
	return usersToDomainModels(pgUsers), nil
}

func (r Repository) UserSave(ctx context.Context, usr User) (*User, error) {
	pgUser := userToPGModel(&usr)

	err := pgUser.Upsert(ctx, r.DB, true, []string{}, boil.Infer(), boil.Blacklist("id"))
	if err != nil {
		return nil, errors.Wrap(err, "upserting user")
	}
	err = pgUser.Reload(ctx, r.DB)
	if err != nil {
		return nil, errors.Wrap(err, "reloading user")
	}
	savedUser := userToDomainModel(pgUser)

	return savedUser, nil
}

func (r Repository) UserDeleteAll(ctx context.Context) (int, error) {
	rowsAffected, err := pgmodel.Users().DeleteAll(ctx, r.DB)
	if err != nil {
		return 0, errors.Wrap(err, "deleting all users")
	}
	return int(rowsAffected), nil
}

func (r Repository) UserDelete(ctx context.Context, id uuid.UUID) (int, error) {
	usr := pgmodel.User{
		ID: id.String(),
	}
	rowsAffected, err := usr.Delete(ctx, r.DB)
	if err != nil {
		return 0, errors.Wrap(err, "deleting user")
	}
	return int(rowsAffected), nil
}

func userToDomainModel(pgUser *pgmodel.User) *User {
	return &User{
		ID:           uuid.MustParse(pgUser.ID),
		Email:        EmailAddress(pgUser.Email),
		PasswordHash: PasswordHash(pgUser.PasswordHash),
	}
}

func userToPGModel(domainUser *User) *pgmodel.User {
	return &pgmodel.User{
		ID:           domainUser.ID.String(),
		Email:        string(domainUser.Email),
		PasswordHash: string(domainUser.PasswordHash),
	}
}

func usersToDomainModels(pgUsers pgmodel.UserSlice) []*User {
	var users []*User
	for _, pgUser := range pgUsers {
		usr := userToDomainModel(pgUser)
		users = append(users, usr)
	}
	return users
}
