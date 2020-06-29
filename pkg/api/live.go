package api

import (
	"net/http"
)


func (s *Server) liveHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("live"))
}



