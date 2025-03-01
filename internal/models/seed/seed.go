package seed

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SeedDatabase(db database.DatabaseManager) {

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
		{ID: utility.GenerateUUID(), Name: "Org1", Email: fmt.Sprintln(utility.RandomString(4) + "@email.com"), Description: "Description1", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org2", Email: fmt.Sprintln(utility.RandomString(4) + "@email.com"), Description: "Description2", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org3", Email: fmt.Sprintln(utility.RandomString(4) + "@email.com"), Description: "Description3", OwnerID: Userid2},
	}

	var existingUser models.User

	query := db.PreloadEntities(nil, &existingUser, "Profile", "Products")
	if err := query.Where("email = ?", user1.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed users
			db.CreateOneRecord(&user1)
			db.CreateOneRecord(&user2)

			// Seed organisations
			for _, org := range organisations {
				db.CreateOneRecord(&org)
			}

			// Seed categories
			for _, category := range categories {
				db.CreateOneRecord(&category)
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

	if err := db.DB().Where("question = ?", faqs[0].Question).First(&models.FAQ{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed faq
			for _, faq := range faqs {
				db.CreateOneRecord(&faq)
			}
		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("FAQ already exist, skipping seeding.")
	}

	//seeding templates
	templates := SeedTemplates()
	if err := db.DB().Where("name = ?", templates[0].Name).First(&models.EmailTemplate{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed templates
			for _, template := range templates {
				db.CreateOneRecord(&template)
			}
		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Templates already exist, skipping seeding.")
	}

}

func SeedTestDatabase(db database.DatabaseManager) {
	roles := []models.Role{
		{ID: int(models.RoleIdentity.User), Name: "user", Description: "user related functions"},
		{ID: int(models.RoleIdentity.SuperAdmin), Name: "super admin", Description: "super admin related functions"},
	}

	var existingRole models.Role
	if err := db.DB().Where("id = ?", roles[0].ID).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			db.CreateMultipleRecords(&roles, len(roles))
		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Roles already exist, skipping seeding.")
	}
}

func SeedOrgRolesAndPermissions(db database.DatabaseManager) {
	var organizations []models.Organisation
	if err := db.DB().Find(&organizations).Error; err != nil {
		fmt.Printf("Error fetching organizations: %v\n", err)
		return
	}

	for _, org := range organizations {
		roles := []models.OrgRole{
			{ID: utility.GenerateUUID(), Name: "Admin", Description: "Administrator Role", OrganisationID: org.ID},
			{ID: utility.GenerateUUID(), Name: "User", Description: "User Role", OrganisationID: org.ID},
		}

		for _, role := range roles {
			isUniqueName, errName := utility.CheckForUniqueness(db, &models.OrgRole{}, "name", role.Name)
			isUniqueId, errId := utility.CheckForUniqueness(db, &models.OrgRole{}, "organisation_id", role.OrganisationID)

			if errName != nil {
				fmt.Printf("Error checking role name uniqueness: %v", errName)
				continue
			}
			if errId != nil {
				fmt.Printf("Error checking role id uniqueness: %v", errId)
				continue
			}
			if isUniqueName && isUniqueId {
				if err := db.CreateOneRecord(&role); err != nil {
					fmt.Printf("Error creating role: %v\n", err)
					continue
				}
			} else {
				fmt.Printf("Role %s already exists in organisation %s\n", role.Name, org.Name)
				var existingRole models.OrgRole
				result := db.DB().Where("name = ? AND organisation_id = ?", role.Name, role.OrganisationID).First(&existingRole)
				if result.Error != nil {
					fmt.Printf("Error fetching existing role: %v\n", result.Error)
					continue
				}
				role.ID = existingRole.ID
				fmt.Printf("Using existing role: ID %s", role.ID)
			}

			permissions := []models.Permission{
				{ID: utility.GenerateUUID(), RoleID: role.ID, Category: "Transactions", PermissionList: map[string]bool{"can_view_transactions": true, "can_edit_transactions": true}},
				{ID: utility.GenerateUUID(), RoleID: role.ID, Category: "Refunds", PermissionList: map[string]bool{"can_view_refunds": true}},
			}

			for _, permission := range permissions {
				isUniqueId, errId := utility.CheckForUniqueness(db, &models.Permission{}, "role_id", permission.RoleID)
				isUniqueCategory, errCategory := utility.CheckForUniqueness(db, &models.Permission{}, "category", permission.Category)
				if errId != nil {
					fmt.Printf("Error checking permissions id uniqueness: %v", errId)
					continue
				}
				if errCategory != nil {
					fmt.Printf("Error checking permissions category uniqueness: %v", errCategory)
					continue
				}
				if isUniqueId && isUniqueCategory {
					if err := db.CreateOneRecord(&permission); err != nil {
						fmt.Printf("Error creating permission: %v\n", err)
					}
				} else {
					fmt.Printf("Category %s already exists for role ID %s in permissions\n", permission.Category, permission.RoleID)
					var existingPermission models.Permission
					result := db.DB().Where("role_id = ? AND category = ?", permission.RoleID, permission.Category).First(&existingPermission)
					if result.Error != nil {
						fmt.Printf("Error fetching existing permission: %v\n", result.Error)
						continue
					}
					permission.ID = existingPermission.ID
					fmt.Printf("Using existing permission: ID %s", permission.ID)
				}
			}
		}
	}
}
