package url

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

type UrlResources struct{}

func (rs UrlResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(middleware.AuthenticateUser).Post("/", rs.createShortURL)
	r.With(middleware.AuthenticateUser).Get("/", rs.getAllShortURLs)
	return r
}

func generateURLForRedirection(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Extract the host (e.g., "localhost:8080")
	host := r.Host

	// Construct the full server URL
	serverURL := fmt.Sprintf("%s://%s/", scheme, host)
	return serverURL
}

func (rs UrlResources) createShortURL(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := types.Url{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if params.URL == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "URL is needed to shorten it...")
		return
	}
	user, ok := middleware.GetUserFromContext(r.Context())
	if ok != true {
		fmt.Println("Unable to get user")
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	params.UserId = user.ID
	baseURL := generateURLForRedirection(r)
	shortenedUrlInfo, err := db.CreateShortUrl(&params, baseURL)
	if err != nil {
		message := err.Error()
		if strings.Contains(err.Error(), "duplicate key error collection:") {
			fmt.Println("Shortened URL already exists for URL:", params.URL)
			message = "Internal Server Error"
		}
		utils.RespondWithError(w, http.StatusInternalServerError, message)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, types.UrlWithShortVersion{
		Url: types.Url{
			ID:     shortenedUrlInfo.Id,
			URL:    params.URL,
			UserId: params.UserId,
		},
		ShortenedUrl: shortenedUrlInfo.ShortenedUrl,
	})
}

func (rs UrlResources) getAllShortURLs(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if ok != true {
		fmt.Println("Unable to get user")
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	baseURL := generateURLForRedirection(r)
	urls, err := db.GetAllShortUrlsForUser(user.ID, baseURL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	fmt.Println(urls)
	utils.RespondWithJSON(w, http.StatusOK, urls)
}
