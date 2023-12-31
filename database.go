package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

// DeleteDB deletes a database if it's found
func DeleteDB(path string) error {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil

}

// CreateUser creates a new chirp and saves it to disk
func (db *DB) CreateUser(email string) (User, error) {
	fmt.Println("creating user")
	db.mux.Lock()
	defer db.mux.Unlock()
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	users, _ := db.GetUsers()
	newUser := User{
		Id:    len(users) + 1,
		Email: email,
	}

	dbData.Users[newUser.Id] = newUser
	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}
	return newUser, nil

}

func (db *DB) GetUsers() ([]User, error) {

	dbData, err := db.loadDB()

	if err != nil {
		return nil, err
	}
	users := make([]User, 0)
	for _, c := range dbData.Users {
		users = append(users, c)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})
	return users, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	fmt.Println("creating chirp")
	db.mux.Lock()
	defer db.mux.Unlock()
	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirpers, _ := db.GetChirps()
	newChirp := Chirp{
		Id:   len(chirpers) + 1,
		Body: body,
	}

	dbData.Chirps[newChirp.Id] = newChirp
	err = db.writeDB(dbData)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil

}

func (db *DB) GetChirpByID(ID int) (Chirp, error) {
	dbData, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}
	val, ok := dbData.Chirps[ID]
	if !ok {
		return Chirp{}, errors.New("not found")
	}
	return val, nil

}

// // GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	dbData, err := db.loadDB()

	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0)
	for _, c := range dbData.Chirps {
		chirps = append(chirps, c)
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	err := os.WriteFile(db.path, []byte("{\"chirps\":{},\"users\":{}}"), 0666)
	if err != nil {
		return err
	}
	return nil
}

// // loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		fmt.Println(err)
	}
	dbData := DBStructure{}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		return DBStructure{}, err
	}
	return dbData, nil
}

// // writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	jsonString, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	writeErr := os.WriteFile(db.path, jsonString, 0666)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
