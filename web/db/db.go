package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/ryomak/tsukemen/web/model"
)

type DBSession struct {
	DB *gorm.DB
}

func NewDBSession() *DBSession {
	db, err := gorm.Open("mysql", "root:@tcp(localhost:3306)/votedb")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.Vote{})

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
	fmt.Println(result)
	return result, nil
}
