package utils

import (
	"net/http"
	"time"
)

func ReadCookie(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func WriteCookie(w http.ResponseWriter, key, value string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Expires:  time.Now().Add(1 * time.Hour), // 1-day expiration
		HttpOnly: true,                          // Prevents access via JavaScript
		Path:     "/",                           // Cookie is valid for the entire site
	}

	// Set the cookie
	http.SetCookie(w, cookie)
}

func DeleteCookie(w http.ResponseWriter, key string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // 1-day expiration
		HttpOnly: true,                           // Prevents access via JavaScript
		Path:     "/",                            // Cookie is valid for the entire site
	}

	// Set the cookie
	http.SetCookie(w, cookie)
}
