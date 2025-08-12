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
	logger.Info(ctx, "responding with data", "data", exampleData)
	json.RespondWithJSON(ctx, w, http.StatusOK, exampleData)
}
