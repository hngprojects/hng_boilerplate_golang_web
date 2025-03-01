package service

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

func CreateTestimonial(db *gorm.DB, req models.TestimonialReq, userId string) (*models.Testimonial, error) {
	testimonial := &models.Testimonial{
		ID:      utility.GenerateUUID(),
		UserID:  userId,
		Name:    req.Name,
		Content: req.Content,
	}

	pdb := inst.InitDB(db)
	err := testimonial.Create(pdb)

	if err != nil {
		return nil, err
	}

	return testimonial, nil

}
