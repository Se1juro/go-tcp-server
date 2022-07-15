package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/upload-files-go/models"
)

type channels struct {
	Channels []int
}

func serverHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("up")
}

func getChannels(w http.ResponseWriter, r *http.Request, server *models.Server) {
	response := channels{Channels: server.GetChannels()}

	json.NewEncoder(w).Encode(response)
}

func Api(server *models.Server) {
	router := mux.NewRouter()
	router.HandleFunc("/api/channels", func(w http.ResponseWriter, r *http.Request) {
		getChannels(w, r, server)
	}).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":3036",
		Handler: router,
	}

	log.Println("Listening on port 3036")
	srv.ListenAndServe()
}
