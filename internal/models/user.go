package models

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
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

func (u *User) AddUserToOrganisation(db *gorm.DB, user interface{}, orgs []interface{}) error {

	// Add user to organisation
	err := db.Model(user).Association("Organisations").Append(orgs...)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByID(db *gorm.DB, userID string) (User, error) {
	var user User

	query := db.Where("id = ?", userID)
	query = postgresql.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) GetUserByEmail(db *gorm.DB, userEmail string) (User, error) {
	var user User

	query := db.Where("email = ?", userEmail)
	query = postgresql.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (u *User) CreateUser(db *gorm.DB) error {

	err := postgresql.CreateOneRecord(db, &u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetSeedUsers(db *gorm.DB) ([]User, error) {
	var users []User

	query := postgresql.PreloadEntities(db, &users, "Profile", "Products", "Organisations")
	query = query.Limit(2)

	if err := query.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (u *User) Update(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &u)
	return err
}

func (u *User) CheckUserIsAdmin(db *gorm.DB) bool {
	return u.Role == int(RoleIdentity.SuperAdmin)
}

func (u *User) GetUserByIDsAdmin(db *gorm.DB, userID, requesterID string) (User, error) {

	var (
		ErrNotFound = errors.New("user not found")
		user        = User{}
	)

	var isOwner bool
	err := db.Model(&Organisation{}).
		Select("count(*) > 0").
		Where("owner_id = ? AND id IN (SELECT organisation_id FROM user_organisations WHERE user_id = ?)", requesterID, userID).
		Find(&isOwner).
		Error
	if err != nil {
		return user, err
	}

	if isOwner {
		query := db.Model(&User{}).
			Joins("INNER JOIN user_organisations uo ON users.id = uo.user_id").
			Where("uo.organisation_id IN (SELECT organisation_id FROM user_organisations WHERE user_id = ?)", userID)
		query = postgresql.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

		if err := query.First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return user, ErrNotFound
			}
			return user, err
		}
		return user, nil
	}

	query := db.Model(&User{}).
		Joins("INNER JOIN user_organisations uo ON users.id = uo.user_id").
		Where("users.id = ? AND users.id = ?", userID, requesterID)
	query = postgresql.PreloadEntities(query, &user, "Profile", "Products", "Organisations")

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *User) DeleteAUser(db *gorm.DB) error {

	err := postgresql.DeleteRecordFromDb(db, u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetProfileID(db *gorm.DB, userID string) (string, error) {
	var user User

	query := db.Where("id = ?", userID)
	query = postgresql.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user.Profile.ID, err
	}

	return user.Profile.ID, nil
}

func (u *User) GetUserWithProfile(db *gorm.DB, userID string) (User, error) {
	var user User

	query := db.Where("id = ?", userID)
	query = postgresql.PreloadEntities(query, &user, "Profile")

	if err := query.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
