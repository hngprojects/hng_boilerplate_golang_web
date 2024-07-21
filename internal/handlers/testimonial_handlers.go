package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joshua468/hng_boilerplate_golang_web/internal/services"
)

// GetTestimonialHandler handles requests to fetch a single testimonial by ID
func GetTestimonialHandler(service services.TestimonialService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": http.StatusUnauthorized,
				"message":     "Unauthorized",
			})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"message":     "Invalid ID format",
			})
			return
		}

		testimonial, err := service.GetTestimonialByID(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"message":     "Internal server error",
			})
			return
		}

		if testimonial == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status_code": http.StatusNotFound,
				"message":     "Testimonial not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusOK,
			"message":     "Testimonial fetched successfully",
			"data": gin.H{
				"id":          testimonial.ID,
				"author":      testimonial.Author,
				"testimonial": testimonial.Testimonial,
				"comments":    testimonial.Comments,
				"created_at":  testimonial.CreatedAt,
			},
		})
	}
}
