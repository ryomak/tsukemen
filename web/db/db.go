package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
  "github.com/ryomak/tsukemen/web/model"
)

type DBSession struct {
	DB *gorm.DB
}

func NewDBSession() *DBSession {
	DBMS := "mysql"
	USER := "root"
	PASS := ""
	PROTOCOL := ""
	DBNAME := "test"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME
	session, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err)
	}
	return &DBSession{DB: session}
}

func (d *DBSession) VoteForCandidate(v model.Vote) error {
	return nil
}
func (d *DBSession) Result() ([]model.Vote, error) {
	return nil, nil
}
