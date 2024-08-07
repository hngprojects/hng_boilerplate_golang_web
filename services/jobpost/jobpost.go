package service

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
)

// type JobPostSummary struct {
// 	ID			string `json:"id"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	Location    string `json:"location"`
// 	Salary      string `json:"salary"`
// }
type JobPostSummary struct {
	ID				 string `json:"id"`
	Title       	 string `json:"title"`
	Location    	 string `json:"location"`
	Description 	 string `json:"description"`
	SalaryRange      string `json:"salary_range"`
}

func CreateJobPost(req models.CreateJobPostModel, db *gorm.DB) (models.JobPost, error) {
	// 	jobpost := models.JobPost{
	// 	ID:          		utility.GenerateUUID(),
	// 	Title:       		req.Title,
	// 	Salary:      		req.Salary,
	// 	JobType:     		req.JobType,
	// 	Location:    		req.Location,
	// 	Deadline:    		req.Deadline,
	//     WorkMode:       	req.WorkMode,
	// 	Experience:			req.Experience,        
	// 	HowToApply:     	req.HowToApply,
	//     JobBenefits:		req.JobBenefits,         
	// 	Description: 		req.Description,
	// 	CompanyName: 		req.CompanyName,
	//     KeyResponsibilities: req.KeyResponsibilities,
	// 	Qualifications:		req.Qualifications,
	// }

	jobpost := models.JobPost{
		ID:          		utility.GenerateUUID(),
		Title:       		req.Title,
		JobMode:       	    req.JobMode,
		JobType:     		req.JobType,
		Location:    		req.Location,
		Deadline:    		req.Deadline,
		SalaryRange:        req.SalaryRange,
		Experience:			req.Experience,                 
		Description: 		req.Description,
		CompanyName: 		req.CompanyName,
	}

	if err := jobpost.CreateJobPost(db); 
	
	err != nil {
		return models.JobPost{}, err
	}

	return jobpost, nil
}

func GetPaginatedJobPosts(c *gin.Context, db *gorm.DB) ([]JobPostSummary, postgresql.PaginationResponse, error) {
	jobpost := models.JobPost{}
	jobPosts, paginationResponse, err := jobpost.FetchAllJobPost(db, c)

	if err != nil {
		return nil, paginationResponse, err
	}

	if len(jobPosts) == 0 {
		return nil, paginationResponse, gorm.ErrRecordNotFound
	}

	var jobPostSummaries []JobPostSummary
	for _, job := range jobPosts {
		// summary := JobPostSummary{
		// 	ID: 		 job.ID,
		// 	Title:       job.Title,
		// 	Description: job.Description,
		// 	Location:    job.Location,
		// 	Salary:      job.Salary,
		// }
		summary := JobPostSummary{
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