package models

import (
	"time"

	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type JobPost struct {
	ID				string 		`gorm:"column:id; type:uuid; primaryKey" json:"job_id`
	Title 			string 		`gorm:"column:title; type:text; not null" json:"title"`
	Description     string 		`gorm:"column:description; type:varchar(255); not null" json:"description"`
	Location        string 		`gorm:"column:location; type:varchar(255); not null" json:"location"`// need review
	Salary 			float64 	`gorm:"column:salary; type:int; not null" json:"salary"`
	JobType 		string 		`gorm:"column:job_type; type:varchar(255); not null" json:"job_type"`
	CompanyName 	string 		`gorm:"column:company_name; type:varchar(255); not null" json:"company_name"`
	CreatedAt       time.Time   `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}


func (j *JobPost) AddJobPost(db *gorm.DB, jobPost JobPost) (JobPost, error) {
	jobPost.ID = utility.GenerateUUID()
	jobPost.CreatedAt = time.Now()

	result := db.Create(&jobPost)
	if result.Error != nil {
		return JobPost{}, result.Error
	}
	return jobPost, nil
}

func (j *JobPost) GetAllJobPosts(db *gorm.DB) ([]JobPost, error) {
	var jobposts []JobPost

	err := db.Find(&jobposts).Error
	if err != nil {
		return nil, err
	}

	return jobposts, nil
}

func (j *JobPost) GetJobPostById(db *gorm.DB, id string) (JobPost, error) {
	var jobpost JobPost

	if err := db.Where("id = ?", id).First(&jobpost).Error; err != nil {
		return jobpost, err
	}

	return jobpost, nil
}

func (j *JobPost) EditJobPostById(db *gorm.DB, id string, updates map[string]interface{}) (JobPost, error) {
	var jobPost JobPost
	if err := db.First(&jobPost, "id = ?", id).Error; err != nil {
		return JobPost{}, err
	}
	if err := db.Model(&jobPost).Updates(updates).Error; err != nil {
		return JobPost{}, err
	}
	jobPost.UpdatedAt = time.Now()
	if err := db.Save(&jobPost).Error; err != nil {
		return JobPost{}, err
	}
	return jobPost, nil
}