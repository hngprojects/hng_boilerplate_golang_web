package job_post

import (
	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"gorm.io/gorm"
)

func CreateJobPost(c *gin.Context, db *gorm.DB) (models.JobPost, error) {
  
}

func FetchAllJobPost(c *gin.Context, db *gorm.DB) (models.JobPost, error) {

}
func FetchJobPostById(c *gin.Context, db *gorm.DB) (models.JobPost, error) {
	id := c.Param("id")
	// 	var user models.User
	var job_post models.JobPost

	if err := db

	// var userResp models.User

	// userResp, err := userResp.GetUserByID(db, userIDStr)
	// if err != nil {
	// 	return userResp, err
	// }

	// return userResp, nil

	// if err := db.Postgresql.First(&user, id).Error; err != nil {
	// 	return models.User{}, errors.New("user not found")
	// }

	// return user, nil
}

func UpdateJobPostById(c *gin.Context, db *gorm.DB) (models.JobPost, error) {
	// id := c.Param("id")
}

// func DeleteJobPostById(c *gin.Context, db *gorm.DB) (models.JobPost, error) {
// 	// id := c.Param("id")
// }
