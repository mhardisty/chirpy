package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type apiConfig struct {
	FileServerHits int
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits++
		fmt.Println(cfg.FileServerHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cfg.AdminPage(w, r)
}

func (cfg *apiConfig) AdminPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
