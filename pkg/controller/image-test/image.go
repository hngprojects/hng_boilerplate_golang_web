package imagetest

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// Retrieve the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// Check if the file type is allowed
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .jpg, .jpeg, and .png files are allowed"})
		return
	}

	// Save the file to a local directory
	savePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": savePath})
}
