package service

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CreateJobPost(req models.CreateJobPostModel, db *gorm.DB) (models.JobPost, error) {
		jobpost := models.JobPost{
		ID:          		utility.GenerateUUID(),
		Title:       		req.Title,
		Salary:      		req.Salary,
		JobType:     		req.JobType,
		Location:    		req.Location,
		Deadline:    		req.Deadline,
	    WorkMode:       	req.WorkMode,
		Experience:			req.Experience,        
		HowToApply:     	req.HowToApply,
	    JobBenefits:		req.JobBenefits,         
		Description: 		req.Description,
		CompanyName: 		req.CompanyName,
	    KeyResponsibilities: req.KeyResponsibilities,
		Qualifications:		req.Qualifications,
	}

	if err := jobpost.CreateJobPost(db); 
	
	err != nil {
		return models.JobPost{}, err
	}

	return jobpost, nil
}