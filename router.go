package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/go-chi/valve"
)

func initHttpRouters() {
	// shutdown signaling.
	valv := valve.New()
	baseCtx := valv.Context()
	_ = baseCtx

	// Multiplexer
	router := chi.NewRouter()

	// Basic settings
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)
	// Basic gracious timing
	router.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})
	router.Use(cors.Handler)

	// Basic Routes
	router.Get("/", IndexPage)

	// Basic Routes Groupings
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/images", ImgApiRoutes())
		r.Mount("/api/credentials", UserApiRoutes())
	})
	log.Println("Starting port", pHttpPort)
	log.Fatal(http.ListenAndServe(":"+pHttpPort, router))
}

func ImgApiRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/upload/{id}", GetOneImage)
	r.Post("/upload", UploadImage)
	r.Get("/", GetAllImages)
	return r
}

func UserApiRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{code}", SetUserCode)
	return r
}
