package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aidenappl/nu-calendar/env"
	"github.com/aidenappl/nu-calendar/middleware"
	"github.com/aidenappl/nu-calendar/routers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Primary Router
	primary := mux.NewRouter()

	// Set used Timezone
	os.Setenv("TZ", "America/New_York")

	// Healthcheck Endpoint
	primary.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Define the API Endpoint
	r := primary.PathPrefix("/nucal/v1.1").Subrouter()

	// Logging Requests
	r.Use(middleware.LoggingMiddleware)

	// Adding Response Headers
	r.Use(middleware.MuxHeaderMiddleware)

	//
	// API Endpoints
	//

	// Add Calendar
	r.HandleFunc("/initializer", routers.HandleInitializer).Methods(http.MethodPost)

	// Get Ref Calendar
	r.HandleFunc("/getCalendar", routers.HandleGetCalendar).Methods(http.MethodGet)

	// ListEvents
	r.HandleFunc("/listEvents", routers.HandleLaunchpad).Methods(http.MethodGet)

	// Edit Reference Events
	r.HandleFunc("/editReferenceEvents", routers.HandleEditReferenceEvents).Methods(http.MethodPost)

	// Retrieve calendar from NUCAL
	r.HandleFunc("/calendar", routers.GetCalendar).Methods(http.MethodGet)

	// Launch Server
	fmt.Printf("✅ APLB NuCalendar API running on port %s\n", env.Port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(originsOk, headersOk, methodsOk)(primary)))
}
