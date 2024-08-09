package service

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)


func CreateJobPost(req models.CreateJobPostModel, db *gorm.DB) (models.JobPost, error) {
	jobpost := models.JobPost{
		ID:          			utility.GenerateUUID(),
		Title:       			req.Title,
		JobMode:       	    	req.JobMode,
		JobType:     			req.JobType,
		Location:    			req.Location,
		Deadline:    			req.Deadline,
		Benefits:				req.Benefits,
		SalaryRange:        	req.SalaryRange,
		Description: 			req.Description,
		CompanyName: 			req.CompanyName,
		ExperienceLevel: 		req.ExperienceLevel,
		KeyResponsibilities: 	req.KeyResponsibilities,
		Qualifications:			req.Qualifications,
	}

	if err := jobpost.CreateJobPost(db); 
	
	err != nil {
		return models.JobPost{}, err
	}

	return jobpost, nil
}

func GetPaginatedJobPosts(c *gin.Context, db *gorm.DB) ([]models.JobPostSummary, postgresql.PaginationResponse, error) {
	jobpost := models.JobPost{}
	jobPosts, paginationResponse, err := jobpost.FetchAllJobPost(db, c)

	if err != nil {
		return nil, paginationResponse, err
	}

	if len(jobPosts) == 0 {
	 	return []models.JobPostSummary{}, paginationResponse, nil
	}

	var jobPostSummaries []models.JobPostSummary
	for _, job := range jobPosts {
		summary := models.JobPostSummary{
			ID: 		 	  job.ID,
			Title:       	  job.Title,
			Description: 	  job.Description,
			Location:    	  job.Location,
			SalaryRange:      job.SalaryRange,
		}
		jobPostSummaries = append(jobPostSummaries, summary)
	}

	return jobPostSummaries, paginationResponse, nil
}

func FetchJobPostByID(db *gorm.DB, id string) (models.JobPost, error) {
	jobpost := models.JobPost{}
	jobpost.ID = id
	err := jobpost.FetchJobPostByID(db)
	if err != nil {
		return models.JobPost{}, err
	}
	return jobpost, nil
}

func UpdateJobPost(db *gorm.DB, jobPost models.JobPost, ID string) (models.JobPost, error) {
	updatedJobPost, err := jobPost.UpdateJobPostByID(db, ID)
	if err != nil {
		return models.JobPost{}, err
	}
	return updatedJobPost, nil
}

func DeleteJobPostByID(db *gorm.DB, ID string) error {
	jobPost := models.JobPost{ID: ID}
	err := jobPost.DeleteJobPostByID(db, ID)
	if err != nil {
		return err
	}
	return nil
}