package app

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func SetUpRouter() {

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/hysterix/test", SimpleHandler).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":5001", router))
}
