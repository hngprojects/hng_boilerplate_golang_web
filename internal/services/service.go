package services

import (
    "context"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/models"
    "gorm.io/gorm"
)

type TestimonialService interface {
    GetTestimonialByID(ctx context.Context, id int) (*models.Testimonial, error)
}

type testimonialService struct {
    db *gorm.DB
}

func NewTestimonialService(db *gorm.DB) TestimonialService {
    return &testimonialService{db: db}
}

func (s *testimonialService) GetTestimonialByID(ctx context.Context, id int) (*models.Testimonial, error) {
    var testimonial models.Testimonial
    if err := s.db.WithContext(ctx).First(&testimonial, id).Error; err != nil {
        return nil, err
    }
    return &testimonial, nil
}
