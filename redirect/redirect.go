package redirect

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/prashant1k99/URL-Shortner/db"
	"github.com/prashant1k99/URL-Shortner/utils"
)

type UserResource struct{}

func (rs UserResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{shortcustURL}", rs.redirectUrl)
	return r
}

func (rs UserResource) redirectUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "shortcustURL")
	if shortUrl == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	urlInfo, err := db.GetUrlFromShortUrl(shortUrl)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	if urlInfo.ShortenedUrl != "" {
		// Track the analytics in a goroutine to avoid blocking the response
		// go trackAnalytics(urlInfo.ID, r)

		http.Redirect(w, r, urlInfo.URL, http.StatusMovedPermanently) // 301 Redirect
		return
	}
}

// getClientIP gets the real client IP address, accounting for possible proxies
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (can contain multiple IPs, comma-separated)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Split the string and take the first part, which is the client's real IP
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Return the full RemoteAddr in case of error
	}

	return ip
}

// func trackAnalytics(urlId primitive.ObjectID, r *http.Request) {
// }
