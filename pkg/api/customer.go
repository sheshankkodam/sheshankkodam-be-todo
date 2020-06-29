package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Customer struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomerResp struct {
	CustomerId string `json:"customer_id"`
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Declare a new Person struct.
	var c Customer

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encryptedPwd := hash(c.Password)
	customerId := c.generateId()

	if insErr := s.db.InsertCustomer(customerId, c.Username, encryptedPwd); insErr != nil {
		http.Error(w, insErr.Error(), http.StatusBadRequest)
		return
	}

	resp, mErr := json.Marshal(CustomerResp{
		CustomerId: customerId,
	})
	if mErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling customer response.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing customer response.")
	}
	return
}

func hash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}


func (c *Customer) generateId() string{
	id := 0

	for _, c := range c.Username {
		id += int(c)
	}

	id += int(':')
	for _, c := range c.Password {
		id += int(c)
	}

	return fmt.Sprint(id)
}


