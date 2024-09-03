package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prashant1k99/URL-Shortner/db"
	"github.com/prashant1k99/URL-Shortner/types"
	"github.com/prashant1k99/URL-Shortner/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.ReadCookie(r, "session_user")
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		coll, err := db.GetCollection("users")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		var user types.User
		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = coll.FindOne(context.TODO(), bson.M{"_id": userObjectId}).Decode(&user)
		if err != nil {
			fmt.Println(err.Error())
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid authentication")
			return
		}
		ctx := context.WithValue(r.Context(), "user", &user)
		// Call the next handler, passing the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NotAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.ReadCookie(r, "session_user")
		if err == nil || userId != "" {
			utils.RespondWithError(w, http.StatusConflict, "User is already authenticated")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value("user").(*types.User)
	return user, ok
}
