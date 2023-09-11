package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	db, dbErr := NewDB(DATABASE_PATH)

	type body struct {
		Body string `json:"email"`
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

	respBody, createErr := db.CreateUser(params.Body)
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
