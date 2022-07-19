package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/upload-files-go/models"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, req)
		})
}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func getChannels(w http.ResponseWriter, r *http.Request, server *models.Server) {
	response := server.GetChannels()

	json.NewEncoder(w).Encode(response)
}

func getHistoryFiles(w http.ResponseWriter, r *http.Request, server *models.Server) {
	response := server.FileHistory
	json.NewEncoder(w).Encode(response)
}

func Api(server *models.Server) {
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/api/channels", func(w http.ResponseWriter, r *http.Request) {
		getChannels(w, r, server)
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/history-files", func(w http.ResponseWriter, r *http.Request) {
		getHistoryFiles(w, r, server)
	}).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":3036",
		Handler: router,
	}
	enableCORS(router)

	log.Println("Listening on port 3036")
	srv.ListenAndServe()
}
