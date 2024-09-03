package user

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prashant1k99/URL-Shortner/utils"
)

type UserResources struct{}

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	APIKey   string `json:"apiKey"`
}

func (rs UserResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-up", rs.signUp)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a user route"))
	})
	return r
}

func (rs UserResources) signUp(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusCreated, user{
		ID:       1,
		Username: "Prashant",
		APIKey:   "some API",
	})
}
