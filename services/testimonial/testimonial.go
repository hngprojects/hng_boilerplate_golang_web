package service

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type TestimonialService interface {
	CreateTestimonial(req models.TestimonialReq, userId string) (*models.Testimonial, error)
}

type testimonialService struct {
	db *gorm.DB
}

func NewTestimonialService(db *gorm.DB) TestimonialService {
	return &testimonialService{db: db}
}

func (s *testimonialService) CreateTestimonial(req models.TestimonialReq, userId string) (*models.Testimonial, error) {
	testimonial := &models.Testimonial{
		ID:      utility.GenerateUUID(),
		UserID:  userId,
		Name:    req.Name,
		Content: req.Content,
	}

	pdb := inst.InitDB(s.db)
	err := testimonial.Create(pdb)

	if err != nil {
		return nil, err
	}

	return testimonial, nil

}
