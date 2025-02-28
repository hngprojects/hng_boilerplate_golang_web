package models

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
)

type User struct {
	ID            string                     `gorm:"type:uuid;primaryKey;unique;not null" json:"id"`
	Name          string                     `gorm:"column:name; type:varchar(255)" json:"name"`
	Email         string                     `gorm:"column:email; type:varchar(255)" json:"email"`
	Password      string                     `gorm:"column:password; type:text; not null" json:"-"`
	Profile       Profile                    `gorm:"foreignKey:Userid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile"`
	Key           Key                        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"key"`
	Region        UserRegionTimezoneLanguage `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"region"`
	DataPrivacy   DataPrivacySettings        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"data_privacy"`
	Organisations []Organisation             `gorm:"many2many:user_organisations;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"organisations" ` // many to many relationship
	Products      []Product                  `gorm:"foreignKey:OwnerID" json:"products"`
	Blogs         []Blog                     `gorm:"foreignKey:AuthorID" json:"blogs"`
	CreatedAt     time.Time                  `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time                  `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
	Role          int                        `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"`
	DeletedAt     gorm.DeletedAt             `gorm:"index" json:"-"`
}

type CreateUserRequestModel struct {
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	UserName    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateUserRequestModel struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	UserName    string `json:"username" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequestModel struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *User) AddUserToOrganisation(db database.DatabaseManager, user interface{}, orgs []interface{}) error {

	// Add user to organisation
	err := db.DB().Model(user).Association("Organisations").Append(orgs...)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByID(db database.DatabaseManager, userID string) (User, error) {
	var user User

	query := db.DB().Where("id = ?", userID)
	query = db.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) GetUserByEmail(db database.DatabaseManager, userEmail string) (User, error) {
	var user User

	query := db.DB().Where("email = ?", userEmail)
	query = db.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) CreateUser(db database.DatabaseManager) error {

	err := db.CreateOneRecord(&u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetSeedUsers(db database.DatabaseManager) ([]User, error) {
	var users []User

	query := db.PreloadEntities(nil, &users, "Profile", "Products", "Organisations")
	query = query.Limit(2)

	if err := query.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (u *User) Update(db database.DatabaseManager) error {
	_, err := db.SaveAllFields(&u)
	return err
}

func (u *User) CheckUserIsAdmin(db database.DatabaseManager) bool {
	return u.Role == int(RoleIdentity.SuperAdmin)
}

func (u *User) GetUserByIDsAdmin(db database.DatabaseManager, userID, requesterID string) (User, error) {

	var (
		ErrNotFound = errors.New("user not found")
		user        = User{}
	)

	var isOwner bool
	err := db.DB().Model(&Organisation{}).
		Select("count(*) > 0").
		Where("owner_id = ? AND id IN (SELECT organisation_id FROM user_organisations WHERE user_id = ?)", requesterID, userID).
		Find(&isOwner).
		Error
	if err != nil {
		return user, err
	}

	if isOwner {
		query := db.DB().Model(&User{}).
			Joins("INNER JOIN user_organisations uo ON users.id = uo.user_id").
			Where("uo.organisation_id IN (SELECT organisation_id FROM user_organisations WHERE user_id = ?)", userID)
		query = db.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

		if err := query.First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return user, ErrNotFound
			}
			return user, err
		}
		return user, nil
	}

	query := db.DB().Model(&User{}).
		Joins("INNER JOIN user_organisations uo ON users.id = uo.user_id").
		Where("users.id = ? AND users.id = ?", userID, requesterID)
	query = db.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *User) DeleteAUser(db database.DatabaseManager) error {

	err := db.DeleteRecordFromDb(u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetProfileID(db database.DatabaseManager, userID string) (string, error) {
	var user User

	query := db.DB().Where("id = ?", userID)
	query = db.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user.Profile.ID, err
	}

	return user.Profile.ID, nil
}

func (u *User) GetUserWithProfile(db database.DatabaseManager, userID string) (User, error) {
	var user User

	query := db.DB().Where("id = ?", userID)
	query = db.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
