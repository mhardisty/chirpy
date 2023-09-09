package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string `json:"body"`
	}

	type success struct {
		Body string `json:"cleaned_body"`
	}

	type error struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respBody := error{
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
		respBody := error{
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

	respBody := success{
		Body: cleanBody(params.Body),
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
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
