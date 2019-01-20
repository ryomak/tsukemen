package db

import (
	"github.com/jinzhu/gorm"
	"github.com/ryomak/tsukemen/web/model"
)

type DBSession struct {
	DB *gorm.DB
}

func NewDBSession() *DBSession {
	db, err := gorm.Open("mysql", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	return &DBSession{DB: db}
}

func (d *DBSession) VoteForCandidate(v model.Vote) error {
	if err := d.DB.Create(&v).Error; err != nil {
		return err
	}
	return nil
}
func (d *DBSession) Result() ([]model.Vote, error) {
	result := make([]model.Vote, 0)
	if err := d.DB.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
