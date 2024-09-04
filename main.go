package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "log"
    "bytes"
)

func sendOTP(phoneNumber string) error {
    apiURL := "https://api.sendpulse.com/sms/send"
    requestBody := map[string]interface{}{
        "phone":   phoneNumber,
        "message": "Your OTP code is 123456",
    }
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return err
    }
    client := &http.Client{}
    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    // Add authentication headers if needed
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // Handle response and errors
    return nil
}

func verifyOTP(phoneNumber string, otp string) bool {
    // Example logic for OTP verification
    // In practice, you'd need to store OTPs and compare
    return otp == "123456" // Dummy check
}

func RegisterPhoneNumber(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        PhoneNumber string `json:"phone_number"`
    }
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    err = sendOTP(requestData.PhoneNumber)
    if err != nil {
        http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
        return
    }

    response := map[string]string{"message": "OTP sent", "status": "success"}
    json.NewEncoder(w).Encode(response)
}

func VerifyPhoneNumber(w http.ResponseWriter, r *http.Request) {
    var requestData struct {
        PhoneNumber string `json:"phone_number"`
        OTP         string `json:"otp"`
    }
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    verified := verifyOTP(requestData.PhoneNumber, requestData.OTP)
    if !verified {
        http.Error(w, "OTP verification failed", http.StatusUnauthorized)
        return
    }

    response := map[string]string{"message": "OTP verified", "status": "success"}
    json.NewEncoder(w).Encode(response)
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/register", RegisterPhoneNumber).Methods("POST")
    router.HandleFunc("/verify", VerifyPhoneNumber).Methods("POST")

    fmt.Println("Server is running on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", router))
}
