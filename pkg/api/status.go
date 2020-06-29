package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Ready status of the service
type readyStatus struct {
	App        string `json:"app"`
	Uptime     string `json:"uptime"`
	Version    string `json:"version"`
}

const AppName = "sheshankkodam-be-todo"

func (s *Server) statusHandler(w http.ResponseWriter, _ *http.Request) {
	statusBytes, err := json.Marshal(readyStatus{
		App:        AppName,
		Uptime:     time.Since(s.startTime).String(),
		Version:    getVersion(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling status response.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(statusBytes); err != nil {
		log.Println("Error writing ready status.")
	}

	return
}

func getVersion() string {
	appVer, readErr := ioutil.ReadFile("VERSION")
	if readErr != nil {
		log.Printf("Unable to read version, error=%s", readErr)
	}
	return string(appVer)
}



