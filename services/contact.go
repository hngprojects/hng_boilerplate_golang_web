package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ContactRequest defines the structure for the contact form request
type ContactRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Message string `json:"message" validate:"required"`
}

// ContactUsHandler handles the contact form submissions
func ContactUsHandler(c *gin.Context, validator *validator.Validate) {
	var req ContactRequest

	// Bind JSON and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "status": http.StatusBadRequest})
		return
	}

	// Validate fields
	// if err := validator.Struct(req); err != nil {
	// 	// Convert the error to a ValidationErrors type
	// 	if validationErrors, ok := err.(validator.ValidationErrors); ok {
	// 		var errors []string
	// 		for _, e := range validationErrors {
	// 			errors = append(errors, e.Error())
	// 		}
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": errors, "status": http.StatusBadRequest})
	// 	} else {
	// 		// If the error is not of type ValidationErrors
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": []string{err.Error()}, "status": http.StatusBadRequest})
	// 	}
	// 	return
	// }

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Inquiry successfully sent", "status": http.StatusOK})
}
