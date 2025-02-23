package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gosb/src/endpoints"
	"gosb/src/global"
	"gosb/src/middlewares"
	"gosb/src/settings"
	"net/http"
	"os"
)

func handleRequests() {
	router := mux.NewRouter()
	//router.StrictSlash()
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Debug level enabled, logging web requests enabled")
		router.Use(middlewares.LoggingMiddleware)
	}

	if settings.EnableCacheHeaders == true {
		router.Use(middlewares.CacheHeadersMiddleware)
	}

	router.HandleFunc("/", indexPage)
	router.HandleFunc(`/api/skipSegments/{shaPrefix:\w{4,32}}`, endpoints.ApiSkipSegmentsEndpoint).Methods(http.MethodGet)

	addrStr := fmt.Sprintf("%s:%d", settings.ListenBind, settings.HttpPort)
	log.Printf("Serving requests on: '%s'", addrStr)
	log.Fatal(http.ListenAndServe(addrStr, router))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "running")
}

func main() {
	{
		var err error
		postgresDSN := os.Getenv("POSTGRES_DSN")
		if postgresDSN == "" {
			log.Fatal("Missing or empty POSTGRES_DSN environment variable")
		}

		global.DB, err = sqlx.Open("postgres", postgresDSN)
		if err != nil {
			log.Fatalf("Failed seeting up db connection: %s", err)
		}

		err = global.DB.Ping()
		if err != nil {
			log.Fatalf("Failed connecting to database: %s", err)
		}

		defer func() {
			err := global.DB.Close()
			if err != nil {
				log.Printf("Error while closing database: %s", err)
			}
		}()
	}
	settings.GetEnvOpts()
	handleRequests()
}
