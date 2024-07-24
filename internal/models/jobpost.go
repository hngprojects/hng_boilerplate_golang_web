package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type JobPost struct {
	ID				string 		`gorm:"column:id; type:uuid; primaryKey" json:"job_id`
	Title 			string 		`gorm:"column:title; type:text; not null" json:"title"`
	Description     string 		`gorm:"column:description; type:varchar(255); not null" json:"description"`
	Location        string 		`gorm:"column:location; type:varchar(255); not null" json:"location"`
	Salary 			float64 	`gorm:"column:salary; type:int; not null" json:"salary"`
	JobType 		string 		`gorm:"column:job_type; type:varchar(255); not null" json:"job_type"`
	CompanyName 	string 		`gorm:"column:company_name; type:varchar(255); not null" json:"company_name"`
	CreatedAt       time.Time   `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

func (j *JobPost) CreateJobPost(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &j)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobPost) FetchAllJobPosts(db *gorm.DB, jobposts *[]JobPost) error {
	err := postgresql.SelectAllFromDb(db, "desc", jobposts, nil)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobPost) FetchJobPost(db *gorm.DB ) error {
	err := postgresql.SelectFirstFromDb(db, &j);

		if err != nil {
		return  err
	}

	return  nil
}
