package jobpost

import (
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CreateJobPost(req models.CreateJobPostModel, db *gorm.DB) (models.JobPost, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	jobpost := models.JobPost{
		ID:                  utility.GenerateUUID(),
		Title:               req.Title,
		JobMode:             req.JobMode,
		JobType:             req.JobType,
		Location:            req.Location,
		Deadline:            req.Deadline,
		Benefits:            req.Benefits,
		SalaryRange:         req.SalaryRange,
		Description:         req.Description,
		CompanyName:         req.CompanyName,
		ExperienceLevel:     req.ExperienceLevel,
		KeyResponsibilities: req.KeyResponsibilities,
		Qualifications:      req.Qualifications,
	}

	if err := jobpost.CreateJobPost(pdb); err != nil {
		return models.JobPost{}, err
	}

	return jobpost, nil
}

func GetPaginatedJobPosts(c *gin.Context, db *gorm.DB) ([]models.JobPostSummary, database.PaginationResponse, error) {
	jobpost := models.JobPost{}
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	jobPosts, paginationResponse, err := jobpost.FetchAllJobPost(pdb, c)

	if err != nil {
		return nil, paginationResponse, err
	}

	if len(jobPosts) == 0 {
		return []models.JobPostSummary{}, paginationResponse, nil
	}

	var jobPostSummaries []models.JobPostSummary
	for _, job := range jobPosts {
		summary := models.JobPostSummary{
			ID:          job.ID,
			Title:       job.Title,
			Description: job.Description,
			Location:    job.Location,
			SalaryRange: job.SalaryRange,
		}
		jobPostSummaries = append(jobPostSummaries, summary)
	}

	return jobPostSummaries, paginationResponse, nil
}

func FetchJobPostByID(db *gorm.DB, id string) (models.JobPost, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	jobpost := models.JobPost{}
	jobpost.ID = id
	err := jobpost.FetchJobPostByID(pdb)
	if err != nil {
		return models.JobPost{}, err
	}
	return jobpost, nil
}

func UpdateJobPost(db *gorm.DB, jobPost models.JobPost, ID string) (models.JobPost, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	updatedJobPost, err := jobPost.UpdateJobPostByID(pdb, ID)
	if err != nil {
		return models.JobPost{}, err
	}
	return updatedJobPost, nil
}

func DeleteJobPostByID(db *gorm.DB, ID string) error {
	jobPost := models.JobPost{ID: ID}
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	err := jobPost.DeleteJobPostByID(pdb, ID)
	if err != nil {
		return err
	}
	return nil
}
