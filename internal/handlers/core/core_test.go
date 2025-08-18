package core_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"discord/internal/db"
	"discord/internal/handlers/core"
	"log/slog"

	"github.com/go-chi/chi/v5"
)

// Mock UserService
type mockUserService struct {
	getUserFunc func(ctx context.Context, id string) (db.User, error)
}

func (m *mockUserService) GetUser(ctx context.Context, id string) (db.User, error) {
	return m.getUserFunc(ctx, id)
}

func TestGetUser_UserFound(t *testing.T) {
	// Arrange
	mockService := &mockUserService{
		getUserFunc: func(ctx context.Context, id string) (db.User, error) {
			return db.User{ID: "1", Email: "test@example.com", Name: "Tester"}, nil
		},
	}

	handler := &core.CoreHandler{
		Log:         slog.Default(),
		Chat:        nil, // not needed for this test
		UserService: mockService,
	}

	r := chi.NewRouter()
	r.Get("/user", handler.GetUser)

	req := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got	db.User 
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal("failed to decode response:", err)
	}

	if got.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", got.Email)
	}
}

// func TestGetExample_Error(t *testing.T) {
// 	mockService := &mockUserService{
// 		getUserFunc: func(ctx context.Context, id string) (db.User, error) {
// 			return db.User{}, errors.New("db error")
// 		},
// 	}
//
// 	handler := &core.CoreHandler{
// 		Log:         slog.Default(),
// 		UserService: mockService,
// 	}
//
// 	r := chi.NewRouter()
// 	r.Get("/user", handler.GetUser)
//
// 	req := httptest.NewRequest("GET", "/user", nil)
// 	w := httptest.NewRecorder()
//
// 	r.ServeHTTP(w, req)
//
// 	if w.Code != http.StatusBadRequest {
// 		t.Fatalf("expected status 400, got %d", w.Code)
// 	}
// }
