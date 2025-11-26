package models

import "time"

type Timeslot struct {
    Time      string `json:"time"`
    Available bool   `json:"available"`
}

type Court struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Type      string `json:"type"`
    Price     int    `json:"price"`
    Available bool   `json:"available"`
}

type Reservation struct {
    ID          string    `json:"id"`
    Date        string    `json:"date"`
    Timeslot    string    `json:"timeslot"`
    CourtID     int       `json:"court_id"`
    CourtName   string    `json:"court_name"`
    CourtType   string    `json:"court_type"`
    UserEmail   string    `json:"user_email"`
    UserName    string    `json:"user_name"`
    UserPhone   string    `json:"user_phone"`
    TotalAmount int       `json:"total_amount"`
    Status      string    `json:"status"` // pending, confirmed, cancelled
    PaymentID   string    `json:"payment_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type PaymentRequest struct {
    Amount      int64  `json:"amount"`
    Currency    string `json:"currency"`
    Description string `json:"description"`
    ReservationID string `json:"reservation_id"`
}

type PaymentResponse struct {
    ClientSecret string `json:"client_secret"`
    PaymentID    string `json:"payment_id"`
}

type CreateReservationRequest struct {
    Date      string `json:"date"`
    Timeslot  string `json:"timeslot"`
    CourtID   int    `json:"court_id"`
    UserEmail string `json:"user_email"`
    UserName  string `json:"user_name"`
    UserPhone string `json:"user_phone"`
}

type ConfirmReservationRequest struct {
    ReservationID string `json:"reservation_id"`
    PaymentID     string `json:"payment_id"`
}