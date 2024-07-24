package service

import (
	"fmt"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func ValidateCreateJobPost(req models.JobPost) error {
	if req.Title == "" || req.Description == "" || req.Location == "" || req.Salary <= 0 || req.JobType == "" || req.CompanyName == "" {
		return fmt.Errorf("all fields are required and must follow their data types")
	}
	return nil
}

func CreateJobPost(req models.JobPost, db *gorm.DB) (models.JobPost, error) {
	jobpost := models.JobPost{
		ID:          utility.GenerateUUID(),
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Salary:      req.Salary,
		JobType:     req.JobType,
		CompanyName: req.CompanyName,
	}
	if err := jobpost.CreateJobPost(db); err != nil {
		return models.JobPost{}, err
	}

	return jobpost, nil
}

func FetchAllJobPost(db *gorm.DB) ([]models.JobPost, error) {
	var jobposts []models.JobPost
	jobpost := models.JobPost{}

	if err := jobpost.FetchAllJobPosts(db, &jobposts); err != nil {
		return nil, err
	}

	return jobposts, nil
}

func FetchJobPostByID(db *gorm.DB, id string) (models.JobPost, error) {
	jobpost := models.JobPost{}
	jobpost.ID = id
	err := jobpost.FetchJobPost(db)
	if err != nil {
		return models.JobPost{}, err
	}
	return jobpost, nil
}




