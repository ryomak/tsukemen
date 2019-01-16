package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/ryomak/tsukemen/web/blockchain"
	"github.com/ryomak/tsukemen/web/db"
	"github.com/ryomak/tsukemen/web/model"
)


type Server struct {
	Store
}

type Store interface {
	VoteForCandidate(model.Vote) error
	Result() ([]model.Vote, error)
}

func (s *Server) VoteForCandidateHandler(w http.ResponseWriter, r *http.Request) {
	vote := new(model.Vote)
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
		db.NewDBSession(),
	}
	blockchainServer := Server{
		blockchain.NewBlockchainSession(),
	}
	r := mux.NewRouter()
	r.HandleFunc("/hello",func(w http.ResponseWriter, r *http.Request) {
		  fmt.Fprintf(w, "Hello, World")
	})
	r.HandleFunc("/db/vote", databaseServer.VoteForCandidateHandler).Methods("POST")
	r.HandleFunc("/db/result", databaseServer.VoteForCandidateHandler).Methods("GET")
	r.HandleFunc("/blockchain/vote", blockchainServer.VoteForCandidateHandler).Methods("POST")
	r.HandleFunc("/blockchain/result", blockchainServer.VoteForCandidateHandler).Methods("GET")
	fmt.Println("run server port:8080")
	http.ListenAndServe(":8080", r)
}
