package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    
    "be-badminton-reservation/handlers"
    "be-badminton-reservation/middleware"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using default values")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    stripeKey := os.Getenv("STRIPE_SECRET_KEY")
    if stripeKey != "" {
        handlers.SetStripeKey(stripeKey)
    }

    r := mux.NewRouter()
    
    api := r.PathPrefix("/api").Subrouter()
    
    api.HandleFunc("/timeslots", handlers.GetTimeslotsHandler).Methods("GET")
    api.HandleFunc("/courts", handlers.GetCourtsHandler).Methods("GET")
    api.HandleFunc("/reservations", handlers.CreateReservationHandler).Methods("POST")
    api.HandleFunc("/reservations", handlers.GetReservationsHandler).Methods("GET")
    api.HandleFunc("/reservations/confirm", handlers.ConfirmReservationHandler).Methods("POST")
    
    api.HandleFunc("/create-payment-intent", handlers.CreatePaymentIntentHandler).Methods("POST")
    api.HandleFunc("/process-dummy-payment", handlers.ProcessDummyPaymentHandler).Methods("POST")
    
    api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status": "ok", "service": "badminton-reservation-api"}`))
    }).Methods("GET")
    
    c := middleware.EnableCORS()
    handler := c.Handler(r)
    
    fmt.Printf(" Badminton Reservation API server running on port %s\n", port)
    fmt.Printf(" Health check: http://localhost%s/api/health\n", port)
    fmt.Printf(" Available endpoints:\n")
    fmt.Printf("   GET  http://localhost%s/api/timeslots?date=YYYY-MM-DD\n", port)
    fmt.Printf("   GET  http://localhost%s/api/courts?date=YYYY-MM-DD&timeslot=HH:MM\n", port)
    fmt.Printf("   POST http://localhost%s/api/reservations\n", port)
    fmt.Printf("   POST http://localhost%s/api/create-payment-intent\n", port)
    
    log.Fatal(http.ListenAndServe(":"+port, handler))
}