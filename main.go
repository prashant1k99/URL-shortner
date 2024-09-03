package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/prashant1k99/URL-Shortner/db"
	"github.com/prashant1k99/URL-Shortner/user"
)

func main() {
	db.ConnectDB()

	router := chi.NewRouter()

	// CORS setup
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	router.Mount("/user", user.UserResources{}.Routes())
	// router.Mount("/shorten-it") // For creating shortened URL
	// router.Mount("/r") // For redirecting user
	defer db.DisconnectDB()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	fmt.Printf("Server is running on port: http://localhost:%v\n", PORT)

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
