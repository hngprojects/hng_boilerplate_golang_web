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

	roles := []models.Role{
		{ID: int(models.RoleIdentity.User), Name: "user", Description: "user related functions"},
		{ID: int(models.RoleIdentity.SuperAdmin), Name: "super admin", Description: "super admin related functions"},
	}

	Userid1 := utility.GenerateUUID()
	user1 := models.User{
		ID:       Userid1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: utility.RandomString(20),
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			AvatarURL: "http://example.com/avatar.jpg",
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product1", Description: "Description1", OwnerID: Userid1},
			{ID: utility.GenerateUUID(), Name: "Product2", Description: "Description2", OwnerID: Userid1},
		},
		Role: int(models.RoleIdentity.User),
	}

	Userid2 := utility.GenerateUUID()
	user2 := models.User{
		ID:       Userid2,
		Name:     "Jane Doe",
		Password: utility.RandomString(20),
		Email:    "jane@example.com",
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "Jane",
			LastName:  "Doe",
			Phone:     "0987654321",
			AvatarURL: "http://example.com/avatar2.jpg",
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product3", Description: "Description3", OwnerID: Userid2},
			{ID: utility.GenerateUUID(), Name: "Product4", Description: "Description4", OwnerID: Userid2},
		},
		Role: int(models.RoleIdentity.SuperAdmin),
	}

	organisations := []models.Organisation{
		{ID: utility.GenerateUUID(), Name: "Org1", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description1", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org2", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description2", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org3", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description3", OwnerID: Userid2},
	}

	var existingRole models.Role
	if err := db.Where("id = ?", roles[0].ID).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			postgresql.CreateMultipleRecords(db, &roles, len(roles))
		} else {
			fmt.Println("An error occurred: ", err)
		}

	} else {
		fmt.Println("Roles already exist, skipping seeding.")
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
			existingUser.AddUserToOrganisation(db, &user1, []interface{}{&organisations[0], &organisations[1]})
			existingUser.AddUserToOrganisation(db, &user2, []interface{}{&organisations[0], &organisations[1], &organisations[2]})
			fmt.Println("Users added to organisations.")

		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Users already exist, skipping seeding.")
	}

}
