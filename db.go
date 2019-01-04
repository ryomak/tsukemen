package main

import "github.com/jinzhu/gorm"

type db gorm.DB

func NewDBServer() *Server {
	return new(Server)
}

func (d *db) VoteForCandidate(v Vote) error {
	return nil
}
func (d *db) Result() []Vote {
	return nil
}
