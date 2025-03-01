package postgresql

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"gorm.io/gorm"
)

var (
	defaultPage  = 1
	defaultLimit = 20
)

func GetPagination(c *gin.Context) database.Pagination {
	var (
		page  *int
		limit *int
	)
	if c.Query("page") != "" {
		pageInt, err := strconv.Atoi(c.Query("page"))
		if err == nil {
			page = &pageInt
		}
	}
	if c.Query("limit") != "" {
		limitInt, err := strconv.Atoi(c.Query("limit"))
		if err == nil {
			limit = &limitInt
		}
	}

	if page != nil && limit != nil {
		return database.Pagination{Page: *page, Limit: *limit}
	} else if page == nil && limit != nil {
		return database.Pagination{Page: defaultPage, Limit: *limit}
	} else if page != nil && limit == nil {
		return database.Pagination{Page: *page, Limit: defaultLimit}
	} else {
		return database.Pagination{Page: defaultPage, Limit: defaultLimit}
	}
}

func (p *Postgresql) SelectAllFromDb(order string, preload string, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	tx := p.DB()
	// apply preloads if needed
	if preload != "" {
		tx = tx.Preload(preload)
	}
	tx = tx.Order("id "+order).Where(query, args...).Find(receiver)
	return tx.Error
}

func (p *Postgresql) SelectAllFromDbWithLimit(order string, limit int, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	tx := p.Db.Order("id "+order).Where(query, args...).Limit(limit).Find(receiver)
	return tx.Error
}

func (p *Postgresql) SelectAllFromDbOrderBy(orderBy, order string, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}
	tx := p.Db.Order(orderBy+" "+order).Where(query, args...).Find(receiver)
	return tx.Error
}

func (p *Postgresql) SelectAllFromByGroup(orderBy, order string, pagination *database.Pagination, receiver interface{}, query interface{}, groupColumn string, args ...interface{}) (database.PaginationResponse, error) {

	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}

	if pagination == nil {
		tx := p.Db.Order(orderBy+" "+order).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
		return database.PaginationResponse{}, tx.Error
	}

	var count int64
	err := p.Db.Model(receiver).Where(query, args...).Group(groupColumn + ", id").Count(&count).Error
	if err != nil {
		return database.PaginationResponse{
			CurrentPage:     pagination.Page,
			PageCount:       pagination.Limit,
			TotalPagesCount: 0,
		}, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	tx := p.Db.Limit(pagination.Limit).Offset((pagination.Page-1)*pagination.Limit).Order(orderBy+" "+order).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
	return database.PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       int(tx.RowsAffected),
		TotalPagesCount: totalPages,
	}, tx.Error
}

func (p *Postgresql) RawSelectAllFromByGroup(orderBy, order string, pagination *database.Pagination, model interface{}, receiver interface{}, groupColumn string, selectQuery string, query string, args ...interface{}) (database.PaginationResponse, error) {

	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}

	if pagination == nil {
		tx := p.Db.Model(model).Order(orderBy+" "+order).Select(selectQuery).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
		return database.PaginationResponse{}, tx.Error
	}

	var count int64
	err := p.Db.Model(model).Where(query, args...).Group(groupColumn + ", id").Count(&count).Error
	if err != nil {
		return database.PaginationResponse{
			CurrentPage:     pagination.Page,
			PageCount:       pagination.Limit,
			TotalPagesCount: 0,
		}, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	tx := p.Db.Model(model).Limit(pagination.Limit).Offset((pagination.Page-1)*pagination.Limit).Order(orderBy+" "+order).Select(selectQuery).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
	return database.PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       int(tx.RowsAffected),
		TotalPagesCount: totalPages,
	}, tx.Error
}

func (p *Postgresql) SelectAllFromDbOrderByPaginated(orderBy, order, filter string, pagination database.Pagination, receiver interface{}, query interface{}, args ...interface{}) (database.PaginationResponse, error) {
	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}
	if pagination.Page <= 0 {
		pagination.Page = defaultPage
	}
	if pagination.Limit < 0 {
		pagination.Limit = defaultLimit
	}

	var count int64

	// optionally apply filter if provided
	db := p.DB()
	if filter != "" {
		db = db.Unscoped().Where(filter)
	}

	err := db.Model(receiver).Where(query, args...).Count(&count).Error
	if err != nil {
		return database.PaginationResponse{
			CurrentPage:     pagination.Page,
			PageCount:       pagination.Limit,
			TotalPagesCount: 0,
		}, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	tx := p.Db.Limit(pagination.Limit).Offset((pagination.Page-1)*pagination.Limit).Order(orderBy+" "+order).Where(query, args...).Find(receiver)
	return database.PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       int(tx.RowsAffected),
		TotalPagesCount: totalPages,
	}, tx.Error
}

func (p *Postgresql) SelectOneFromDb(receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := p.Db.Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}

func (p *Postgresql) SelectLatestFromDb(receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := p.Db.Order("id desc").Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}

func (p *Postgresql) SelectRandomFromDb(receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := p.Db.Order("rand()").Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}

func (p *Postgresql) SelectFirstFromDb(receiver interface{}) error {
	tx := p.Db.First(receiver)
	return tx.Error
}

func (p *Postgresql) CheckExists(receiver interface{}, query interface{}, args ...interface{}) bool {

	tx := p.Db.Where(query, args...).First(receiver)
	return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
}

func (p *Postgresql) CheckExistsInTable1(table string, query interface{}, args ...interface{}) bool {
	var result interface{}
	tx := p.Db.Table(table).Where(query, args...).Take(&result)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return false
		} else {
			fmt.Println("tx error", tx.Error.Error())
		}
	}

	return true
}

func (p *Postgresql) CheckExistsInTable(table string, query interface{}, args ...interface{}) bool {
	var result map[string]interface{}
	tx := p.Db.Table(table).Where(query, args...).Take(&result)
	return tx.RowsAffected != 0
}

func (p *Postgresql) PreloadEntities(db *gorm.DB, model interface{}, preloads ...string) *gorm.DB {
	if db == nil {
		db = p.Db
	}
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}
