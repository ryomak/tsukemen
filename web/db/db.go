package db

import (
	"github.com/gocql/gocql"
	"github.com/ryomak/tsukemen/web/model"
)

type DBSession struct {
	DB *gocql.Session
}

func NewDBSession() *DBSession {
	cluster := gocql.NewCluster("localhost:7000", "localhost:7001")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	return &DBSession{DB: session}
}

func (d *DBSession) VoteForCandidate(v model.Vote) error {
	if err := d.DB.Query(`INSERT INTO tweet (user, candidate_id) VALUES (?, ?)`, v.User, v.CandidateID).Exec(); err != nil {
		return err
	}
	return nil
}
func (d *DBSession) Result() ([]model.Vote, error) {
	result := make([]model.Vote, 0)
	vote := model.Vote{}
	iter := d.DB.Query(`SELECT user,cadidate_id FROM votes`).Iter()
	for iter.Scan(&vote.User, &vote.CandidateID) {
		result = append(result, vote)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return result, nil
}
