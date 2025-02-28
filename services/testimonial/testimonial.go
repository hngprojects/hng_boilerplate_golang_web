package service

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type TestimonialService interface {
	CreateTestimonial(req models.TestimonialReq, userId string) (*models.Testimonial, error)
}

type testimonialService struct {
	db database.DatabaseManager
}

func NewTestimonialService(db database.DatabaseManager) TestimonialService {
	return &testimonialService{db: db}
}

func (t *testimonialService) CreateTestimonial(req models.TestimonialReq, userId string) (*models.Testimonial, error) {
	testimonial := &models.Testimonial{
		ID:      utility.GenerateUUID(),
		UserID:  userId,
		Name:    req.Name,
		Content: req.Content,
	}

	err := testimonial.Create(t.db)

	if err != nil {
		return nil, err
	}

	return testimonial, nil

}
