package models

import (
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type JobPostSummary struct {
	ID				 string `json:"id"`
	Title       	 string `json:"title"`
	Location    	 string `json:"location"`
	Description 	 string `json:"description"`
	SalaryRange      string `json:"salary_range"`
}

type JobPost struct {
	ID                  string    `gorm:"type:uuid; primaryKey" json:"id"`
	Title               string    `gorm:"column:title; type:varchar(255); not null" json:"title"`
	SalaryRange         string    `gorm:"column:salary_range; type:varchar(255); null" json:"salary_range"`
	JobType             string    `gorm:"column:job_type; type:varchar(50); not null" json:"job_type"`
	Location            string    `gorm:"column:location; type:varchar(255); not null" json:"location"`
	Deadline            time.Time `gorm:"column:deadline; not null" json:"deadline"`
	JobMode             string    `gorm:"column:job_mode; type:varchar(50); null" json:"job_mode"`
	WorkMode            string    `gorm:"column:work_mode; type:varchar(50); null" json:"work_mode"`
	CompanyName         string    `gorm:"column:company_name; type:varchar(255); not null" json:"company_name"`
	Description         string    `gorm:"column:description; type:varchar(500); not null" json:"description"`
	Benefits            string    `gorm:"column:benefits; type:varchar(500)" json:"benefits"`
	ExperienceLevel     string    `gorm:"column:experience; type:varchar(50); null" json:"experience_level"`
	KeyResponsibilities string    `gorm:"column:key_responsibilities; type:varchar(500)" json:"key_responsibilities"`
	Qualifications      string    `gorm:"column:qualifications; type:varchar(500)" json:"qualifications"`
	CreatedAt           time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type CreateJobPostModel struct {
	Title               string    `json:"title" validate:"required,min=2,max=255"`
	SalaryRange         string    `json:"salary_range" validate:"required"`
	JobType             string    `json:"job_type" validate:"required"`
	Location            string    `json:"location" validate:"required,min=2,max=255"`
	Deadline            time.Time `json:"deadline" validate:"required"`
	JobMode             string    `json:"job_mode" validate:"required"`
	Benefits            string    `json:"benefits" validate:"required,min=2,max=500"`
	CompanyName         string    `json:"company_name" validate:"required,min=2,max=255"`
	Description         string    `json:"description" validate:"required,min=2,max=500"`
	ExperienceLevel     string    `json:"experience_level" validate:"required,min=2,max=500"`
	KeyResponsibilities string    `json:"key_responsibilities" validate:"required,min=2,max=500"`
	Qualifications      string    `json:"qualifications" validate:"required,min=2,max=500"`
}

func (m *CreateJobPostModel) Sanitize() {
	v := reflect.ValueOf(m).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			cleanedValue := utility.CleanStringInput(field.String())
			field.SetString(cleanedValue)
		}
	}
}

func (j *JobPost) CreateJobPost(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &j)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobPost) FetchAllJobPost(db *gorm.DB, c *gin.Context) ([]JobPost, postgresql.PaginationResponse, error) {
	var jobPosts []JobPost

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(
		db,
		"created_at",
		"desc",
		pagination,
		&jobPosts,
		nil,
	)

	if err != nil {
		return nil, paginationResponse, err
	}

	return jobPosts, paginationResponse, nil
}

func (j *JobPost) FetchJobPostByID(db *gorm.DB) error {
	err := postgresql.SelectFirstFromDb(db, &j)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobPost) UpdateJobPostByID(db *gorm.DB, ID string) (JobPost, error) {
	j.ID = ID

	exists := postgresql.CheckExists(db, &JobPost{}, "id = ?", ID)
	if !exists {
		return JobPost{}, gorm.ErrRecordNotFound
	}

	_, err := postgresql.SaveAllFields(db, j)
	if err != nil {
		return JobPost{}, err
	}

	updatedJobPost := JobPost{}
	err = db.First(&updatedJobPost, "id = ?", ID).Error
	if err != nil {
		return JobPost{}, err
	}
	return updatedJobPost, nil
}

func (j *JobPost) DeleteJobPostByID(db *gorm.DB, ID string) error {

	exists := postgresql.CheckExists(db, &JobPost{}, "id = ?", ID)
	if !exists {
		return gorm.ErrRecordNotFound
	}

	err := postgresql.DeleteRecordFromDb(db, &j)

	if err != nil {
		return err
	}

	return nil
}
