package service

import (
	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type TestimonialService interface {
	GetUserTestimonials(userID string) ([]models.Testimonial, error)
}

type TestimonialServiceImpl struct {
	CreateTestimonial(req models.TestimonialReq, userId string) (*models.Testimonial, error)
}

type testimonialService struct {
	db *gorm.DB
}

func NewTestimonialService(db *gorm.DB) TestimonialService {
	return &TestimonialServiceImpl{db: db}
}

func CreateTestimonial(db *gorm.DB, req models.TestimonialReq, userId string) (*models.Testimonial, error) {
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

func (s *TestimonialServiceImpl) GetUserTestimonials(userID string) ([]models.Testimonial, error) {
	var testimonials []models.Testimonial

	if err := s.db.Where("user_id = ?", userID).Find(&testimonials).Error; err != nil {
		return nil, err
	}

	return testimonials, nil
}

