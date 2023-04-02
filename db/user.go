package db

import (
	"time"

	"github.com/gocql/gocql"
)

type UserEntity struct {
	Id        string    `cql:"id"`
	Name      string    `cql:"name"`
	Age       int       `cql:"age"`
	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

type UserRepository struct {
	session *gocql.Session
}

func NewUserRepository(session *gocql.Session) UserRepository {
	return UserRepository{
		session: session,
	}
}

func (r UserRepository) FetchUsers(count int) ([]UserEntity, error) {
	query := r.session.Query("SELECT id, name, age, created_at, updated_at FROM app.users LIMIT ?", count)
	iter := query.Iter()

	var users []UserEntity
	var user UserEntity
	for iter.Scan(
		&user.Id,
		&user.Name,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	) {
		users = append(users, user)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r UserRepository) FindUser(userId string) (*UserEntity, error) {
	query := r.session.Query("SELECT id, name, age, created_at, updated_at FROM app.users WHERE id = ?", userId)
	iter := query.Iter()

	var user UserEntity
	if iter.Scan(
		&user.Id,
		&user.Name,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	) {
		return &user, nil
	} else if err := iter.Close(); err != nil {
		return nil, err
	}

	return nil, gocql.ErrNotFound
}
