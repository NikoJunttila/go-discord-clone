package user

import (
	"context"
	"discord/internal/db"
	"fmt"
)

// User is your domain entity
type User struct {
	ID    int64
	Email string
	Name  string
}

// Interface for easier testing
type Service interface {
	GetUser(ctx context.Context, id string) (db.User, error)
}

// Concrete implementation
type UserService struct {
	queries Queries // wrap your db.Queries here
}

type Queries interface {
	GetUser(ctx context.Context, id string) (db.User, error)
}

func NewUserService(q Queries) *UserService {
	return &UserService{queries: q}
}

func (s *UserService) GetUser(ctx context.Context, id string) (db.User, error) {
	fmt.Println("here I do stuff then return user and error")
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
