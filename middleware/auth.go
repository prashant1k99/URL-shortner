package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prashant1k99/URL-Shortner/db"
	"github.com/prashant1k99/URL-Shortner/types"
	"github.com/prashant1k99/URL-Shortner/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.ReadCookie(r, "session_user")
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userObjectId, err := primitive.ObjectIDFromHex(userId)

		user, err := db.GetUserById(userObjectId)
		if err != nil {
			fmt.Println(err.Error())
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid authentication")
			return
		}
		userInfo := types.User{
			ID:       user.ID,
			Username: user.Username,
		}
		ctx := context.WithValue(r.Context(), "user", userInfo)
		// Call the next handler, passing the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsNotAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.ReadCookie(r, "session_user")
		if err == nil || userId != "" {
			utils.RespondWithError(w, http.StatusConflict, "User is already authenticated")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func IsAuthenticed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.ReadCookie(r, "session_user")
		if err != nil || userId == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "User is not authenticated")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value("user").(types.User)
	return &user, ok
}
