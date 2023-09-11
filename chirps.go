package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func postChirpHandler(w http.ResponseWriter, r *http.Request) {
	db, dbErr := NewDB(DATABASE_PATH)

	type body struct {
		Body string `json:"body"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	if dbErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		respBody := errorResponse{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respBody := errorResponse{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		respBody := errorResponse{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}

	respBody, createErr := db.CreateChirp(cleanBody(params.Body))
	fmt.Println("response created")
	if createErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("marshaling")

	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}

func getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("handling")
	db, err := NewDB(DATABASE_PATH)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ID := chi.URLParam(r, "id")
	fmt.Println(ID)
	IDint, err := strconv.Atoi(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chirp, getErr := db.GetChirpByID(IDint)
	if getErr != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("Marshalling!!!!")
	marshalled, mErr := json.Marshal(chirp)
	if mErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshalled)

}

func getChirpHandler(w http.ResponseWriter, r *http.Request) {
	db, err := NewDB(DATABASE_PATH)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chirps, chirpErr := db.GetChirps()
	if chirpErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshalled, mErr := json.Marshal(chirps)
	if mErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshalled)
}

func cleanBody(body string) string {
	println(body)
	bodyWords := strings.Split(body, " ")
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range bodyWords {
		for _, badWord := range profaneWords {
			if strings.ToLower(word) == badWord {
				bodyWords[i] = "****"
			}
		}
	}
	ret := strings.Join(bodyWords, " ")
	println(ret)
	return ret
}
