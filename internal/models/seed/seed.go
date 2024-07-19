package seed

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SeedDatabase(db *gorm.DB) {
	// instantiate uuid

	Userid1 := utility.GenerateUUID()
	user1 := models.User{
		Userid: Userid1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			AvatarURL: "http://example.com/avatar.jpg",
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product1", Description: "Description1", Userid: Userid1},
			{ID: utility.GenerateUUID(), Name: "Product2", Description: "Description2", Userid: Userid1},
		},
	}

	Userid2 := utility.GenerateUUID()
	user2 := models.User{
		Userid: Userid2,
		Name:   "Jane Doe",
		Email:  "jane@example.com",
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "Jane",
			LastName:  "Doe",
			Phone:     "0987654321",
			AvatarURL: "http://example.com/avatar2.jpg",
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product3", Description: "Description3", Userid: Userid2},
			{ID: utility.GenerateUUID(), Name: "Product4", Description: "Description4", Userid: Userid2},
		},
	}

	organisations := []models.Organisation{
		{Orgid: utility.GenerateUUID(), Name: "Org1", Description: "Description1"},
		{Orgid: utility.GenerateUUID(), Name: "Org2", Description: "Description2"},
		{Orgid: utility.GenerateUUID(), Name: "Org3", Description: "Description3"},
	}

	var existingUser models.User
	if err := db.Preload("Profile").Preload("Products").Where("email = ?", user1.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			postgresql.CreateOneRecord(db, &user1)
			postgresql.CreateOneRecord(db, &user2)
			for _, org := range organisations {
				postgresql.CreateOneRecord(db, &org)
			}
			fmt.Println("Users and organisations seeded.")

			// Add users to organisations

			// add user1 to two organization
			postgresql.AddUserToOrganisation(db, organisations[0].Orgid, user1.Userid)
			postgresql.AddUserToOrganisation(db, organisations[1].Orgid, user1.Userid)

			// Add user2 to the three organization
			postgresql.AddUserToOrganisation(db, organisations[0].Orgid, user2.Userid)
			postgresql.AddUserToOrganisation(db, organisations[1].Orgid, user2.Userid)
			postgresql.AddUserToOrganisation(db, organisations[2].Orgid, user2.Userid)
			fmt.Println("Users added to organisations.")

		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Users already exist, skipping seeding.")
	}

}
