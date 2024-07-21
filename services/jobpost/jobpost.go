package jobpost

import (

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"gorm.io/gorm"
)

func CreateJobPost(db *gorm.DB, reqBody models.JobPost) (models.JobPost, error) {
	jobPost, err := reqBody.AddJobPost(db, reqBody)
	if err != nil {
		return jobPost, err
	}

	return jobPost, nil
}

func FetchAllJobPost(db *gorm.DB) ([]models.JobPost, error) {
  var jobPost models.JobPost
  jobPosts, err := jobPost.GetAllJobPosts(db) 
  if err != nil {
    return nil, err
  }

  return jobPosts, nil
}

func FetchJobPostById(db *gorm.DB, reqBody string) (models.JobPost, error) {
  var jobPost models.JobPost

  jobPost, err := jobPost.GetJobPostById(db, reqBody)
  if err != nil {
    return jobPost, err
  }

  return jobPost, nil
}

func UpdateJobPostById(db *gorm.DB, id string, updates map[string]interface{}) (models.JobPost, error) {
	var jobPost models.JobPost

	jobPost, err := jobPost.EditJobPostById(db, id, updates)
	if err != nil {
		return jobPost, err
	}

	return jobPost, nil
}