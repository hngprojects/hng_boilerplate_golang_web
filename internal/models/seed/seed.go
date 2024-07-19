package seed

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
)

func SeedDatabase(db *gorm.DB) {
	// Check and seed users
	userId1 := uuid.New().String()
	user1 := models.User{
		UserID: userId1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Profile: models.Profile{
			ID:        uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			AvatarURL: "http://example.com/avatar.jpg",
		},
		Products: []models.Product{
			{ID: uuid.New().String(), Name: "Product1", Description: "Description1", UserID: userId1},
			{ID: uuid.New().String(), Name: "Product2", Description: "Description2", UserID: userId1},
		},
	}

	userId2 := uuid.New().String()
	user2 := models.User{
		UserID: userId2,
		Name:   "Jane Doe",
		Email:  "jane@example.com",
		Profile: models.Profile{
			ID:        uuid.New().String(),
			FirstName: "Jane",
			LastName:  "Doe",
			Phone:     "0987654321",
			AvatarURL: "http://example.com/avatar2.jpg",
		},
		Products: []models.Product{
			{ID: uuid.New().String(), Name: "Product3", Description: "Description3", UserID: userId2},
			{ID: uuid.New().String(), Name: "Product4", Description: "Description4", UserID: userId2},
		},
	}

	organisations := []models.Organisation{
		{OrgID: uuid.New().String(), Name: "Org1", Description: "Description1"},
		{OrgID: uuid.New().String(), Name: "Org2", Description: "Description2"},
		{OrgID: uuid.New().String(), Name: "Org3", Description: "Description3"},
	}

	var existingUser models.User
	if err := db.Preload("Profile").Preload("Products").Where("email = ?", user1.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			db.Create(&user1)
			db.Create(&user2)
			for _, org := range organisations {
				db.Create(&org)
			}
			fmt.Println("Users and organisations seeded.")

			// Add users to organisations

			// add user1 to two organization
			models.AddUserToOrganisation(db, organisations[0].OrgID, user1.UserID)
			models.AddUserToOrganisation(db, organisations[1].OrgID, user1.UserID)

			// Add user2 to the three organization
			models.AddUserToOrganisation(db, organisations[0].OrgID, user2.UserID)
			models.AddUserToOrganisation(db, organisations[1].OrgID, user2.UserID)
			models.AddUserToOrganisation(db, organisations[2].OrgID, user2.UserID)
			fmt.Println("Users added to organisations.")

		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Users already exist, skipping seeding.")
	}
}
