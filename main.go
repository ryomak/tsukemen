package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Vote struct {
	ID          uint
	CandidateID uint
}

type Server struct {
	Store
}

type Store interface {
	VoteForCandidate(Vote) error
	Result() ([]Vote, error)
}

func (s *Server) VoteForCandidateHandler(w http.ResponseWriter, r *http.Request) {
	vote := new(Vote)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// anything to db
	if err := s.VoteForCandidate(*vote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ResultHandler(w http.ResponseWriter, r *http.Request) {
	votes, err := s.Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(votes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func main() {
	databaseServer := Server{
		NewDBSession(),
	}
	blockchainServer := Server{
		NewBlockchainSession(),
	}
	r := mux.NewRouter()
	r.HandleFunc("/db/vote", databaseServer.VoteForCandidateHandler).Methods("POST")
	r.HandleFunc("/db/result", databaseServer.VoteForCandidateHandler).Methods("GET")
	r.HandleFunc("/blockchain/vote", blockchainServer.VoteForCandidateHandler).Methods("POST")
	r.HandleFunc("/blockchain/result", blockchainServer.VoteForCandidateHandler).Methods("GET")
}
