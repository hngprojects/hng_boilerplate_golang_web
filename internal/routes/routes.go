package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/handlers"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/middleware"
    "github.com/joshua468/hng_boilerplate_golang_web/internal/services"
    "gorm.io/gorm"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
    testimonialService := services.NewTestimonialService(db)
    router.GET("/api/v1/testimonials/:id", middleware.AuthMiddleware(), handlers.GetTestimonialHandler(testimonialService))
}
