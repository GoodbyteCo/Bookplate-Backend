package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-chi/jwtauth"
	"github.com/holopollock/Bookplate/Middleware"
	"github.com/holopollock/Bookplate/routes"
	"github.com/holopollock/Bookplate/routes/auth"
	"github.com/holopollock/Bookplate/utils"
	"net/http"
	"time"
)

func init() {
	utils.Migrate()
}

func main() {
	r := chi.NewRouter()
	c := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(c.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Group(func(r chi.Router) {

		r.Get("/auth", auth.Auth)
		r.Get("/auth/callback", auth.AuthCallback)

		r.Post("/add/book", routes.AddBook)
	})

	r.Route("/book/{bookID}", func(r chi.Router) {
		r.Use(Middleware.ArticleCtx)
		r.Get("/", routes.GetBook)

	})

	_ = http.ListenAndServe(":8080", r)
}