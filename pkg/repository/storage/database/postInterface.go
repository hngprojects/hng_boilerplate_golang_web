package database

import (
	"gorm.io/gorm"
)

// brought here to avoid cycli
type Pagination struct {
	Page  int
	Limit int
}
type PaginationResponse struct {
	CurrentPage     int `json:"current_page"`
	PageCount       int `json:"page_count"`
	TotalPagesCount int `json:"total_pages_count"`
}

// Reader interface for read operation interface
type Reader interface {
	SelectAllFromDb(order string, preload string, receiver interface{}, query interface{}, args ...interface{}) error
	SelectAllFromDbWithLimit(order string, limit int, receiver interface{}, query interface{}, args ...interface{}) error
	SelectAllFromDbOrderBy(orderBy, order string, receiver interface{}, query interface{}, args ...interface{}) error
	SelectAllFromDbOrderByPaginated(orderBy, order, filter string, pagination Pagination, receiver interface{}, query interface{}, args ...interface{}) (PaginationResponse, error)
	SelectAllFromByGroup(orderBy, order string, pagination *Pagination, receiver interface{}, query interface{}, groupColumn string, args ...interface{}) (PaginationResponse, error)
	RawSelectAllFromByGroup(orderBy, order string, pagination *Pagination, model interface{}, receiver interface{}, groupColumn string, selectQuery string, query string, args ...interface{}) (PaginationResponse, error)
	SelectOneFromDb(receiver interface{}, query interface{}, args ...interface{}) (error, error)
	SelectLatestFromDb(receiver interface{}, query interface{}, args ...interface{}) (error, error)
	SelectRandomFromDb(receiver interface{}, query interface{}, args ...interface{}) (error, error)
	SelectFirstFromDb(receiver interface{}) error
}

// Preloader interface for preloading operations
type Preloader interface {
	PreloadEntities(db *gorm.DB, model interface{}, preloads ...string) *gorm.DB
}

// Writer interface for write operations
type Writer interface {
	// Create operations
	CreateOneRecord(model interface{}) error
	CreateMultipleRecords(models interface{}, length int) error

	// Update operations
	SaveAllFields(model interface{}) (*gorm.DB, error)
	UpdateFields(model interface{}, updates interface{}, id string) (*gorm.DB, error)
	SaveAllModelsFields(models []interface{}) (*gorm.DB, error)

	// Delete operations
	DeleteRecordFromDb(record interface{}) error
	HardDeleteRecordFromDb(record interface{}) error
}

// Counter interface for database counting operations
type Counter interface {
	CheckExists(receiver interface{}, query interface{}, args ...interface{}) bool
	CheckExistsInTable(table string, query interface{}, args ...interface{}) bool
	CountRecords(model interface{}) (int64, error)
	CountSpecificRecords(model interface{}, query interface{}) (int64, error)
}

// DBAccessor interface for getting the database connection
type DBAccessor interface {
	DB() *gorm.DB
}

type DatabaseManager interface {
	Reader
	Writer
	Counter
	Preloader
	DBAccessor
}
