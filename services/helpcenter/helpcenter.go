package service

import (

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
	"gorm.io/gorm"
)

type HelpCntSummary struct {
	ID			string `json:"id"`
	Title       string `json:"title"`
	Content 	string `json:"content"`
	Author      string `json:"author"`
}

func CreateHelpCenterTopic(req models.CreateHelpCenter, db *gorm.DB) (models.HelpCenter, error) {
	helpCnt := models.HelpCenter{
		ID:          		utility.GenerateUUID(),
		Title:       		req.Title,	
		Content:       		req.Content,	
		Author:       		req.Author,	
	}

	if err := helpCnt.CreateHelpCenterTopic(db);

	err != nil {
		return models.HelpCenter{}, err
	}

	return helpCnt, nil
}

func GetPaginatedTopics(c *gin.Context, db *gorm.DB) ([]HelpCntSummary, postgresql.PaginationResponse, error) {
	helpCnt := models.HelpCenter{}
	helpCnts, paginationResponse, err := helpCnt.FetchAllTopics(db, c)

	if err != nil {
		return nil, paginationResponse, err
	}

	if len(helpCnts) == 0 {
		return nil, paginationResponse, gorm.ErrRecordNotFound
	}
	
	var topicSummaries []HelpCntSummary
	for _, Hlp := range helpCnts {
		summary := HelpCntSummary{
			ID: 		 Hlp.ID,
			Title:       Hlp.Title,
			Content:     Hlp.Content,
			Author:      Hlp.Author,
		}
		topicSummaries = append(topicSummaries, summary)
	}

	return topicSummaries, paginationResponse, nil
}

func FetchTopicByID(db *gorm.DB, id string) (models.HelpCenter, error) {
	helpCnt := models.HelpCenter{}
	helpCnt.ID = id
	err := helpCnt.FetchTopicByID(db)
	if err != nil {
		return models.HelpCenter{}, err
	}
	return helpCnt, nil
}

func SearchHelpCenterTopics(c *gin.Context, db *gorm.DB, query string) ([]HelpCntSummary, postgresql.PaginationResponse, error) {
	var helpCnt models.HelpCenter
	topics, paginationResponse, err := helpCnt.SearchHelpCenterTopics(db, c, query)

	if err != nil {
		return nil, paginationResponse, err
	}

	if len(topics) == 0 {
		return nil, paginationResponse, gorm.ErrRecordNotFound
	}

	var topicSummaries []HelpCntSummary
	for _, topic := range topics {
		summary := HelpCntSummary{
			ID:      topic.ID,
			Title:   topic.Title,
			Content: topic.Content,
			Author:  topic.Author,
		}
		topicSummaries = append(topicSummaries, summary)
	}

	return topicSummaries, paginationResponse, nil
}

func UpdateTopic(db *gorm.DB, helpCnt models.HelpCenter, ID string) (models.HelpCenter, error) {
	updatedHelpCnt, err := helpCnt.UpdateTopicByID(db, ID)
	if err != nil {
		return models.HelpCenter{}, err
	}
	return updatedHelpCnt, nil
}

func DeleteTopicByID(db *gorm.DB, ID string) error {
	helpCnt := models.HelpCenter{ID: ID}
	err := helpCnt.DeleteTopicByID(db, ID)
	if err != nil {
		return err
	}
	return nil
}