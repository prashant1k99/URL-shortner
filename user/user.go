package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/prashant1k99/URL-Shortner/db"
	"github.com/prashant1k99/URL-Shortner/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResources struct{}

type user struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}

type userWithPassword struct {
	user
	Password string `json:"password"`
}

type userWithAPI struct {
	user
	APIKey string `json:"apiKey"`
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
	decoder := json.NewDecoder(r.Body)

	params := userWithPassword{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if params.Username == "" || params.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Username and password are required to Signup")
		return
	}
	userCollection, err := db.GetCollection("users")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	encryptedPass, err := utils.EncryptPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	result, err := userCollection.InsertOne(context.TODO(), bson.M{
		"username": params.Username,
		"password": encryptedPass,
	})
	if err != nil {
		message := err.Error()
		if strings.Contains(err.Error(), "duplicate key error collection:") {
			message = "User with this username already exists"
		}
		utils.RespondWithError(w, http.StatusInternalServerError, message)
		return
	}
	fmt.Printf("Inserted a document: %v\n", result.InsertedID)
	utils.RespondWithJSON(w, http.StatusCreated, user{
		ID:       result.InsertedID.(primitive.ObjectID),
		Username: params.Username,
	})
}
