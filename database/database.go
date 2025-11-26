package database

import (
    "sync"
    "strconv"
    
    "be-badminton-reservation/models"
)

var (
    reservations []models.Reservation
    mutex        sync.RWMutex
)

// sample courts
var Courts = []models.Court{
    {ID: 1, Name: "Court 1 - Basic", Type: "basic", Price: 20},
    {ID: 2, Name: "Court 2 - Basic", Type: "basic", Price: 20},
    {ID: 3, Name: "Court 3 - VIP", Type: "vip", Price: 35},
    {ID: 4, Name: "Court 4 - VIP", Type: "vip", Price: 35},
    {ID: 5, Name: "Court 5 - VVIP", Type: "vvip", Price: 50},
}

func CreateReservation(reservation models.Reservation) {
    mutex.Lock()
    defer mutex.Unlock()
    reservations = append(reservations, reservation)
}

func GetReservations() []models.Reservation {
    mutex.RLock()
    defer mutex.RUnlock()
    return reservations
}

func GetReservationByID(id string) *models.Reservation {
    mutex.RLock()
    defer mutex.RUnlock()
    for _, r := range reservations {
        if r.ID == id {
            return &r
        }
    }
    return nil
}

func UpdateReservationStatus(id string, status string, paymentID string) bool {
    mutex.Lock()
    defer mutex.Unlock()
    for i := range reservations {
        if reservations[i].ID == id {
            reservations[i].Status = status
            reservations[i].PaymentID = paymentID
            return true
        }
    }
    return false
}

func IsCourtAvailable(date string, timeslot string, courtID int) bool {
    mutex.RLock()
    defer mutex.RUnlock()
    
    for _, r := range reservations {
        if r.Date == date && r.Timeslot == timeslot && r.CourtID == courtID && r.Status == "confirmed" {
            return false
        }
    }
    return true
}

func GetReservationsByDateAndTime(date string, timeslot string) []models.Reservation {
    mutex.RLock()
    defer mutex.RUnlock()
    
    var result []models.Reservation
    for _, r := range reservations {
        if r.Date == date && r.Timeslot == timeslot && r.Status == "confirmed" {
            result = append(result, r)
        }
    }
    return result
}

func GenerateTimeslots(date string) []models.Timeslot {
    timeslots := []models.Timeslot{}
    
    // Generate timeslots from 8 AM to 10 PM
    for hour := 8; hour < 22; hour++ {
        timeStr := formatTime(hour)
        available := isTimeslotAvailable(date, timeStr)
        timeslots = append(timeslots, models.Timeslot{
            Time:      timeStr,
            Available: available,
        })
    }
    
    return timeslots
}

func GetAvailableCourts(date string, timeslot string) []models.Court {
    availableCourts := []models.Court{}
    
    for _, court := range Courts {
        courtCopy := court
        courtCopy.Available = IsCourtAvailable(date, timeslot, court.ID)
        availableCourts = append(availableCourts, courtCopy)
    }
    
    return availableCourts
}

func isTimeslotAvailable(date string, timeslot string) bool {
    reservations := GetReservationsByDateAndTime(date, timeslot)
    // Assume maximum 2 reservations per timeslot for demo
    return len(reservations) < 2
}

func formatTime(hour int) string {
    if hour < 10 {
        return "0" + strconv.Itoa(hour) + ":00"
    }
    return strconv.Itoa(hour) + ":00"
}

func GetCourtByID(id int) *models.Court {
    for _, court := range Courts {
        if court.ID == id {
            return &court
        }
    }
    return nil
}