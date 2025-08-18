package user

import (
	"context"
	"discord/internal/db"
	"log/slog"
)

type User struct {
	ID    int
	Email string
	Name  string
}
type UserService struct {
	DB *db.Queries
}

func NewUserService(db *db.Queries) *UserService {
	return &UserService{DB: db}
}
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	slog.Info("Calling getUser", "id", id)
	return &User{
		ID:    id,
		Email: "xdd",
		Name:  "xdd2",
	}, nil
}
