package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/stripe/stripe-go/v78"
    "github.com/stripe/stripe-go/v78/paymentintent"
    
    "be-badminton-reservation/models"
)

var stripeSecretKey string

func SetStripeKey(key string) {
    stripeSecretKey = key
}

func CreatePaymentIntentHandler(w http.ResponseWriter, r *http.Request) {
    var req models.PaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // For demo purposes, if Stripe key is not set, use dummy payment
    if stripeSecretKey == ""{
        createDummyPayment(w, req)
        return
    }
    
    // Real Stripe payment
    stripe.Key = stripeSecretKey
    
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(req.Amount),
        Currency: stripe.String(req.Currency),
        Description: stripe.String(req.Description),
        AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
            Enabled: stripe.Bool(true),
        },
    }
    
    pi, err := paymentintent.New(params)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := models.PaymentResponse{
        ClientSecret: pi.ClientSecret,
        PaymentID:    pi.ID,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func createDummyPayment(w http.ResponseWriter, req models.PaymentRequest) {
    // Simulate payment processing
    paymentID := "pay_dummy_" + strconv.FormatInt(req.Amount, 10) + "_" + strconv.Itoa(int(req.Amount))
    
    response := models.PaymentResponse{
        ClientSecret: "dummy_client_secret_" + paymentID,
        PaymentID:    paymentID,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func ProcessDummyPaymentHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ReservationID string `json:"reservation_id"`
        Amount        int64  `json:"amount"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Simulate payment processing delay
    paymentID := "dummy_pay_" + strconv.FormatInt(req.Amount, 10) + "_" + strconv.Itoa(int(req.Amount))
    
    response := map[string]interface{}{
        "success":      true,
        "payment_id":   paymentID,
        "reservation_id": req.ReservationID,
        "message":      "Payment processed successfully",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}