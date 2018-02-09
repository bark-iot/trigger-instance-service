package main

import (
	"database/sql"
	"net/http"
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@db:5432/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":80", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/trigger-instances", a.createTriggerInstance).Methods("POST")
}

func (a *App) createTriggerInstance(w http.ResponseWriter, r *http.Request) {
	var t TriggerInstance
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := t.createTriggerInstance(a.DB); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, _ := json.Marshal(t)
	log.Printf("Created Trigger Instance %s", response)
	respondWithJSON(w, http.StatusCreated, t)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
