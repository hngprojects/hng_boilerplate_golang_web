package models

import (
	"time"

	"gorm.io/gorm"
)

type JobPost struct {
	ID				string 		`gorm:"column:id; type:uuid; primaryKey" json:"job_id`
	Title 			string 		`gorm:"column:title; type:text; not null" json:"title"`
	Description     string 		`gorm:"column:description; type:varchar(255); not null" json:"description"`
	Location        string 		`gorm:"column:location; type:varchar(255); not null" json:"location"`// need review
	Salary 			string 		`gorm:"column:salary; type:number; not null" json:"salary"`
	JobType 		string 		`gorm:"column:job_title; type:varchar(255); not null" json:"job_title"`
	CompanyName 	string 		`gorm:"column:company_name; type:varchar(255); not null" json:"company_name"`
	CreatedAt       time.Time   `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}


func (u *JobPost) CreateJobPost(db *gorm.DB, error{
	err := db.Model().Association("User").Append(user.user_id)
	if err != nil {
		return err
	}
	return nil
})
func (u *JobPost) FetchAllJobPost(db *storage.Database, error){
		err := db.Postgresql.Find(&JobPost).Error; err != nil {
		return models.JobPost{}, errors.New("No job post found!")
	}
}

func (u *JobPost) FetchJobPostById(db *storage.Database, error){
		err := db.Postgresql.First(&JobPost, id).Error; err != nil {
		return models.JobPost{}, errors.New("Job post not found!")
	}
}

func (u *JobPost) UpdateJobPostById(db *storage.Database, error{
		err := db.Postgresql(&JobPost).Updates(JobPost{Name: "hello", Age: 18, Active: false})
})




// func (u *User) GetUserByID(db *gorm.DB, userID string) (User, error) {
	// 	var user User
	
	// 	if err := db.Preload("Profile").Preload("Products").Preload("Organisations").Where("id = ?", userID).First(&user).Error; err != nil {
		// 		return user, err
		// 	}
		
		// 	return user, nil
		// }

		// func (u *User) AddUserToOrganisation(db *gorm.DB, user interface{}, orgs []interface{}) error {
		
		// 	// Add user to organisation
		// 	err := db.Model(user).Association("Organisations").Append(orgs...)
		// 	if err != nil {
		// 		return err
		// 	}
		
		// 	return nil
		// }