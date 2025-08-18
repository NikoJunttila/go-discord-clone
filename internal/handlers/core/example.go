package core

import (
	"discord/internal/handlers/json"
	logger "discord/pkg/logging"
	"errors"
	"math/rand"
	"net/http"
)

func (c *CoreHandler) GetExample(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	num := rand.Intn(10)
	if num > 5 {
		err := errors.New("test error")
		json.RespondWithError(ctx, w, http.StatusBadRequest, "example error", err)
		return
	}
	exampleData := struct {
		Title  string
		Number int
	}{
		Title:  "example json",
		Number: num,
	}
	user, err := c.UserService.GetUser(ctx, "1")
	if err != nil {
		json.RespondWithError(ctx, w, http.StatusBadRequest, "error finding user", err)
		return
	}
	logger.InfoCTX(ctx, "responding with data", "data", exampleData, "user", user)
	json.RespondWithJSON(ctx, w, http.StatusOK, user)
}

func (c *CoreHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	//get id from body or json request
	user, err := c.UserService.GetUser(r.Context(), "1")
	if err != nil {
		json.RespondWithError(r.Context(), w, http.StatusBadRequest, "Failed to find user", err)
		return
	}
	json.RespondWithJSON(r.Context(), w, http.StatusOK, user)
}
