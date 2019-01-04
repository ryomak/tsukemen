package main

import (
	"github.com/gorilla/mux"
)

type Vote struct {
	ID          uint
	CandidateID uint
}

type Server struct {
	Store Store
}

type Store interface {
	VoteForCandidate(Vote) error
	Result() []Vote
}

func main() {

	db := NewDBServer()
	r := mux.NewRouter()
	r.HandlerFunc("/"func()error)
}
