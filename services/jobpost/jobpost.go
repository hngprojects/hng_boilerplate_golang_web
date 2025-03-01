package jobpost

import (
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

// JobPostService defines the interface for job post operations
type JobPostService interface {
	CreateJobPost(req models.CreateJobPostModel) (models.JobPost, error)
	GetPaginatedJobPosts(c *gin.Context) ([]models.JobPostSummary, database.PaginationResponse, error)
	FetchJobPostByID(id string) (models.JobPost, error)
	UpdateJobPost(jobPost models.JobPost, ID string) (models.JobPost, error)
	DeleteJobPostByID(ID string) error
}

// jobPostService is the concrete implementation of JobPostService
type jobPostService struct {
	db *gorm.DB
}

// NewJobPostService creates a new instance of jobPostService
func NewJobPostService(db *gorm.DB) JobPostService {
	return &jobPostService{db: db}
}

func (s *jobPostService) CreateJobPost(req models.CreateJobPostModel) (models.JobPost, error) {
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

	pdb := inst.InitDB(s.db)
	if err := jobpost.CreateJobPost(pdb); err != nil {
		return models.JobPost{}, err
	}

	return jobpost, nil
}

func (s *jobPostService) GetPaginatedJobPosts(c *gin.Context) ([]models.JobPostSummary, database.PaginationResponse, error) {
	jobpost := models.JobPost{}
	pdb := inst.InitDB(s.db)
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

func (s *jobPostService) FetchJobPostByID(id string) (models.JobPost, error) {
	jobpost := models.JobPost{ID: id}
	pdb := inst.InitDB(s.db)
	if err := jobpost.FetchJobPostByID(pdb); err != nil {
		return models.JobPost{}, err
	}
	return jobpost, nil
}

func (s *jobPostService) UpdateJobPost(jobPost models.JobPost, ID string) (models.JobPost, error) {
	pdb := inst.InitDB(s.db)
	updatedJobPost, err := jobPost.UpdateJobPostByID(pdb, ID)
	if err != nil {
		return models.JobPost{}, err
	}
	return updatedJobPost, nil
}

func (s *jobPostService) DeleteJobPostByID(ID string) error {
	jobPost := models.JobPost{ID: ID}
	pdb := inst.InitDB(s.db)
	return jobPost.DeleteJobPostByID(pdb, ID)
}
