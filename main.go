package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
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
	apiRouter.Post("/validate_chirp", validateChirpHandler)
	adminRouter.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.ServeHTTP(w, r)
	})

	corsMux := middlewareCors(r)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux}
	server.ListenAndServe()
}
