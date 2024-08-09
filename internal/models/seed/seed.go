package seed

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SeedDatabase(db *gorm.DB) {

	SeedTestDatabase(db)

	Userid1 := utility.GenerateUUID()
	Userid2 := utility.GenerateUUID()

	categories := []models.Category{
		{ID: utility.GenerateUUID(), Name: "Fashion"},
		{ID: utility.GenerateUUID(), Name: "Groceries"},
		{ID: utility.GenerateUUID(), Name: "Appliances"},
	}

	// Create users
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
		Region: models.UserRegionTimezoneLanguage{
			ID:         utility.GenerateUUID(),
			UserID:     Userid1,
			RegionID:   utility.GenerateUUID(),
			LanguageID: utility.GenerateUUID(),
			TimezoneID: utility.GenerateUUID(),
		},
		DataPrivacy: models.DataPrivacySettings{
			ID:     utility.GenerateUUID(),
			UserID: Userid1,
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product1", Description: "Description1", Price: 45.33, OwnerID: Userid1},
			{ID: utility.GenerateUUID(), Name: "Product2", Description: "Description2", Price: 45.33, OwnerID: Userid1},
		},
		Role: int(models.RoleIdentity.User),
	}

	user2 := models.User{
		ID:       Userid2,
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: utility.RandomString(20),
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "Jane",
			LastName:  "Doe",
			Phone:     "0987654321",
			AvatarURL: "http://example.com/avatar2.jpg",
		},
		DataPrivacy: models.DataPrivacySettings{
			ID:     utility.GenerateUUID(),
			UserID: Userid1,
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product3", Description: "Description3", Price: 45.33, OwnerID: Userid2},
			{ID: utility.GenerateUUID(), Name: "Product4", Description: "Description4", Price: 45.33, OwnerID: Userid2},
		},
		Role: int(models.RoleIdentity.SuperAdmin),
	}

	// Create organisations and categories
	organisations := []models.Organisation{
		{ID: utility.GenerateUUID(), Name: "Org1", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description1", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org2", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description2", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org3", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description3", OwnerID: Userid2},
	}

	var existingUser models.User

	query := postgresql.PreloadEntities(db, &existingUser, "Profile", "Products")
	if err := query.Where("email = ?", user1.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed users
			postgresql.CreateOneRecord(db, &user1)
			postgresql.CreateOneRecord(db, &user2)

			// Seed organisations
			for _, org := range organisations {
				postgresql.CreateOneRecord(db, &org)
			}

			// Seed categories
			for _, category := range categories {
				postgresql.CreateOneRecord(db, &category)
			}

			fmt.Println("Users, organisations and categories seeded.")

			// // Add users to organisations
			existingUser.AddUserToOrganisation(db, &user1, []interface{}{&organisations[0], &organisations[1]})
			existingUser.AddUserToOrganisation(db, &user2, []interface{}{&organisations[0], &organisations[1]})
			fmt.Println("Users added to organisations.")

			// Add products to categories
			for _, product := range user1.Products {
				product.AddProductToCategory(db, []interface{}{&categories[0], &categories[1]})
			}

			for _, product := range user2.Products {
				product.AddProductToCategory(db, []interface{}{&categories[0], &categories[1]})
			}

			SeedOrgRolesAndPermissions(db)

			fmt.Println("Products added to categories.")
		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Users already exist, skipping seeding.")
	}

	faqs := []models.FAQ{
		{ID: utility.GenerateUUID(), Question: "What is the latest fashion trend?", Answer: "The latest fashion trend is...", Category: "Policies"},
		{ID: utility.GenerateUUID(), Question: "What are the best grocery stores?", Answer: "The best grocery stores are...", Category: "Policies"},
		{ID: utility.GenerateUUID(), Question: "How do I choose the right appliance?", Answer: "To choose the right appliance, you should...", Category: "Policies"},
	}

	if err := db.Where("question = ?", faqs[0].Question).First(&models.FAQ{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed faq
			for _, faq := range faqs {
				postgresql.CreateOneRecord(db, &faq)
			}
		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("FAQ already exist, skipping seeding.")
	}

	//seeding templates
	templates := SeedTemplates()
	if err := db.Where("name = ?", templates[0].Name).First(&models.EmailTemplate{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed templates
			for _, template := range templates {
				postgresql.CreateOneRecord(db, &template)
			}
		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Templates already exist, skipping seeding.")
	}

}

func SeedTestDatabase(db *gorm.DB) {

	roles := []models.Role{
		{ID: int(models.RoleIdentity.User), Name: "user", Description: "user related functions"},
		{ID: int(models.RoleIdentity.SuperAdmin), Name: "super admin", Description: "super admin related functions"},
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

}

func SeedOrgRolesAndPermissions(db *gorm.DB) {

	var organizations []models.Organisation
	if err := db.Find(&organizations).Error; err != nil {
		fmt.Printf("Error fetching organizations: %v\n", err)
		return
	}

	for _, org := range organizations {
		roles := []models.OrgRole{
			{ID: utility.GenerateUUID(), Name: "Admin", Description: "Administrator Role", OrganisationID: org.ID},
			{ID: utility.GenerateUUID(), Name: "User", Description: "User Role", OrganisationID: org.ID},
		}

		for _, role := range roles {
			if err := postgresql.CreateOneRecord(db, &role); err != nil {
				fmt.Printf("Error creating role: %v\n", err)
				continue
			}

			permissions := []models.Permission{
				{ID: utility.GenerateUUID(), RoleID: role.ID, Category: "Transactions", PermissionList: map[string]bool{"can_view_transactions": true, "can_edit_transactions": true}},
				{ID: utility.GenerateUUID(), RoleID: role.ID, Category: "Refunds", PermissionList: map[string]bool{"can_view_refunds": true}},
			}

			for _, permission := range permissions {
				if err := postgresql.CreateOneRecord(db, &permission); err != nil {
					fmt.Printf("Error creating permission: %v\n", err)
				}
			}
		}
	}
}
