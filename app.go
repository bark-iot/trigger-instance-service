package main

import (
	"app/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"database/sql"
	"net/http"
	"encoding/json"
	"log"
	"os"
	"strings"
	"crypto/hmac"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"time"
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

func (a *App) authorize(r *http.Request, bodyBytes []byte) error {
	auth := r.Header.Get("Authorization")
	log.Printf("Attempt to authorize message %s with header %s", bodyBytes, auth)
	keys := strings.Split(strings.Replace(auth, "Bark ", "", 1), ":")
	if len(keys) < 2 {
		return errors.New("wrong Authorization header")
	}

	// Find House
	h := models.House{Key: keys[0]}
	if err := h.GetHouse(a.DB); err != nil {
		return err
	}

	// Check hmac signature
	mac := hmac.New(sha256.New, []byte(h.Secret))
	mac.Write(bodyBytes)
	expectedMAC := mac.Sum(nil)
	if strings.Compare(keys[1], hex.EncodeToString(expectedMAC)) != 0 {
		return errors.New("wrong signature")
	}

	return nil
}

func (a *App) createTriggerInstance(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	// Authorize
	if err := a.authorize(r, bodyBytes); err != nil {
		log.Printf("Not authorized with error: %v", err)
		respondWithError(w, http.StatusForbidden, "Invalid Authorization header")
		return
	}
	defer r.Body.Close()

	var ti models.TriggerInstance
	var dat map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &dat); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	ti.InputData = dat["input_data"].(string)

	// Get trigger ID
	d := models.Device{Token: dat["device_token"].(string)}
	if err := d.GetDevice(a.DB); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid device token")
		return
	}
	t := models.Trigger{DeviceID: d.ID, Key: dat["key"].(string)}
	if err := t.GetTrigger(a.DB); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid trigger key")
		return
	}
	ti.TriggerID = t.ID
	ti.CreatedAt = time.Now()

	if err := ti.CreateTriggerInstance(a.DB); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response, _ := json.Marshal(ti)
	log.Printf("Created Trigger Instance %s", response)
	respondWithJSON(w, http.StatusCreated, ti)
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
