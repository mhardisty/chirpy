package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const DATABASE_PATH string = "database.json"

func main() {
	DeleteDB(DATABASE_PATH)
	db, err := NewDB(DATABASE_PATH)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbFileErr := db.ensureDB()
	if dbFileErr != nil {
		log.Fatal(dbFileErr)
		return
	}
	my_err := db.ensureDB()
	if my_err != nil {
		log.Fatal(my_err)
	}

	apiCfg := apiConfig{
		FileServerHits: 0,
	}

	r := chi.NewRouter()

	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	r.Handle("/app/*", apiCfg.middlewareMetrics(appHandler))
	r.Handle("/app", apiCfg.middlewareMetrics(appHandler))

	r.Mount("/api/", apiRouter)
	r.Mount("/admin/", adminRouter)

	apiRouter.Get("/healthz", readinessHandler)
	apiRouter.Get("/chirps", getChirpHandler)
	apiRouter.Get("/chirps/{id}", getChirpByIDHandler)
	apiRouter.Post("/chirps", postChirpHandler)
	adminRouter.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.ServeHTTP(w, r)
	})

	corsMux := middlewareCors(r)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux}
	server.ListenAndServe()
}
