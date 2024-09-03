package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/prashant1k99/URL-Shortner/db"
	"github.com/prashant1k99/URL-Shortner/middleware"
	"github.com/prashant1k99/URL-Shortner/types"
	"github.com/prashant1k99/URL-Shortner/utils"
)

type UserResources struct{}

func (rs UserResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(middleware.IsNotAuthenticated).Post("/sign-up", rs.signUp)
	r.With(middleware.IsNotAuthenticated).Post("/login", rs.logIn)
	r.With(middleware.IsAuthenticed).Delete("/sign-out", rs.signOut)

	r.With(middleware.AuthenticateUser).Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUserFromContext(r.Context())
		if ok != true {
			fmt.Println("Unable to get user")
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, user)
	})
	return r
}

func (rs UserResources) signUp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := types.UserWithPassword{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if params.Username == "" || params.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Username and password are required to Signup")
		return
	}
	userId, err := db.CreateUser(&params)
	if err != nil {
		message := err.Error()
		if strings.Contains(err.Error(), "duplicate key error collection:") {
			message = "User with this username already exists"
		}
		utils.RespondWithError(w, http.StatusInternalServerError, message)
		return
	}
	utils.WriteCookie(w, "session_user", userId.Hex())
	utils.RespondWithJSON(w, http.StatusCreated, types.User{
		ID:       userId,
		Username: params.Username,
	})
}

func (rs UserResources) signOut(w http.ResponseWriter, _ *http.Request) {
	utils.DeleteCookie(w, "session_user")
	type res struct {
		Status string `json:"status"`
	}
	utils.RespondWithJSON(w, http.StatusOK, res{
		Status: "Ok",
	})
}

func (rs UserResources) logIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := types.UserWithPassword{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if params.Username == "" || params.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Username and password are required to Signup")
		return
	}

	user, err := db.GetUserByUsername(params.Username)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	isPasswordValid, err := utils.ComparePassword(params.Password, user.Password)
	if err != nil {
		fmt.Println("Error occured while comparing pass", err)
	}
	if !isPasswordValid {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid username and password")
		return
	}

	utils.WriteCookie(w, "session_user", user.ID.Hex())
	utils.RespondWithJSON(w, http.StatusCreated, types.User{
		ID:       user.ID,
		Username: params.Username,
	})
}
