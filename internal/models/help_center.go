package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)


type HelpCenter struct {
	ID        string    `gorm:"type:uuid; primaryKey" json:"id"`
	Title     string    `gorm:"column:title; type:varchar(255); not null" json:"title"`
	Content   string    `gorm:"column:content; type:varchar(255); not null" json:"content"`
	Author    string    `gorm:"column:author; type:varchar(255); not null" json:"author"`
	CreatedAt time.Time `gorm:"column:created_at; not null; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; null; autoUpdateTime" json:"updated_at"`
}

type CreateHelpCenter struct {
	Title string `json:"title" validate:"required,min=2,max=255"`
	Content string `json:"content" validate:"required,min=2,max=255"`
	Author string `json:"author" validate:"required,min=2,max=255"`
}

func (j *HelpCenter) CreateHelpCenterTopic(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &j)

	if err != nil {
		return err
	}

	return nil
}

func (j *HelpCenter) FetchAllTopics(db *gorm.DB, c *gin.Context) ([]HelpCenter, postgresql.PaginationResponse, error) {
	var helpCntTopics []HelpCenter

	pagination := postgresql.GetPagination(c)

	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(
		db,
		"created_at", 
		"desc",       
		pagination,   
		&helpCntTopics,  
		nil,          
	)

	if err != nil {
		return nil, paginationResponse, err
	}
	
	return helpCntTopics, paginationResponse, nil
}

func (j *HelpCenter) FetchTopicByID(db *gorm.DB ) error {
		err := postgresql.SelectFirstFromDb(db, &j);

		if err != nil {
		return  err
	}

	return  nil
}

func (h *HelpCenter) SearchHelpCenterTopics(db *gorm.DB, c *gin.Context, query string) ([]HelpCenter, postgresql.PaginationResponse, error) {
	var helpCntTopics []HelpCenter
	pagination := postgresql.GetPagination(c)
	searchQuery := "%" + query + "%"
	whereClause := "title ILIKE ?"
	paginationResponse, err := postgresql.SelectAllFromDbOrderByPaginated(
		db,
		"created_at",
		"desc",
		pagination,
		&helpCntTopics,
		whereClause,
		searchQuery,
	)
	if err != nil {
		return nil, paginationResponse, err
	}
	
	return helpCntTopics, paginationResponse, nil
}

func (j *HelpCenter) UpdateTopicByID(db *gorm.DB, ID string) (HelpCenter, error) {
	j.ID = ID 
	
    exists := postgresql.CheckExists(db, &HelpCenter{}, "id = ?", ID)
    if !exists {
        return HelpCenter{}, gorm.ErrRecordNotFound
    }

    _, err := postgresql.SaveAllFields(db, j)
    if err != nil {
        return HelpCenter{}, err
    }

    updatedHelpCenter := HelpCenter{}
    err = db.First(&updatedHelpCenter, "id = ?", ID).Error
    if err != nil {
        return HelpCenter{}, err
    }

    return updatedHelpCenter, nil
}

func (j *HelpCenter) DeleteTopicByID(db *gorm.DB, ID string) error {

	exists := postgresql.CheckExists(db, &HelpCenter{}, "id = ?", ID)
	if !exists {
		return gorm.ErrRecordNotFound
	}

	err := postgresql.DeleteRecordFromDb(db, &j)

	if err != nil {
		return err
	}

	return nil
}
