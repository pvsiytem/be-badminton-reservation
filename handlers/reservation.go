package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "strconv"
    
    "be-badminton-reservation/database"
    "be-badminton-reservation/models"
)

func GetTimeslotsHandler(w http.ResponseWriter, r *http.Request) {
    date := r.URL.Query().Get("date")
    if date == "" {
        http.Error(w, "Date parameter is required", http.StatusBadRequest)
        return
    }
    
    timeslots := database.GenerateTimeslots(date)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(timeslots)
}

func GetCourtsHandler(w http.ResponseWriter, r *http.Request) {
    date := r.URL.Query().Get("date")
    timeslot := r.URL.Query().Get("timeslot")
    
    if date == "" || timeslot == "" {
        http.Error(w, "Date and timeslot parameters are required", http.StatusBadRequest)
        return
    }
    
    courts := database.GetAvailableCourts(date, timeslot)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courts)
}

func CreateReservationHandler(w http.ResponseWriter, r *http.Request) {
    var req models.CreateReservationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // validate required fields
    if req.Date == "" || req.Timeslot == "" || req.CourtID == 0 || req.UserEmail == "" || req.UserName == "" {
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }
    
    if !database.IsCourtAvailable(req.Date, req.Timeslot, req.CourtID) {
        http.Error(w, "Court is no longer available", http.StatusConflict)
        return
    }
    
    // court details
    court := database.GetCourtByID(req.CourtID)
    if court == nil {
        http.Error(w, "Invalid court ID", http.StatusBadRequest)
        return
    }
    
    // create reservation with pending status
    reservation := models.Reservation{
        ID:          generateReservationID(),
        Date:        req.Date,
        Timeslot:    req.Timeslot,
        CourtID:     req.CourtID,
        CourtName:   court.Name,
        CourtType:   court.Type,
        UserEmail:   req.UserEmail,
        UserName:    req.UserName,
        UserPhone:   req.UserPhone,
        TotalAmount: court.Price,
        Status:      "pending", 
        CreatedAt:   time.Now(),
    }
    
    database.CreateReservation(reservation)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reservation)
}

func ConfirmReservationHandler(w http.ResponseWriter, r *http.Request) {
    var req models.ConfirmReservationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    if req.ReservationID == "" || req.PaymentID == "" {
        http.Error(w, "Reservation ID and Payment ID are required", http.StatusBadRequest)
        return
    }
    
    success := database.UpdateReservationStatus(req.ReservationID, "confirmed", req.PaymentID)
    if !success {
        http.Error(w, "Reservation not found", http.StatusNotFound)
        return
    }
    
    reservation := database.GetReservationByID(req.ReservationID)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reservation)
}

func GetReservationsHandler(w http.ResponseWriter, r *http.Request) {
    reservations := database.GetReservations()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reservations)
}

func generateReservationID() string {
    return "RES" + strconv.FormatInt(time.Now().UnixNano(), 10)
}