package billing

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/inst"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/database"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

// BillingService defines the interface for billing operations.
type BillingService interface {
	CreateBilling(req models.CreateBillingRequest, userId string) (models.BillingResponse, error)
	DeleteBilling(BillingId string, userId string) error
	GetBillings(c *gin.Context) (int, database.PaginationResponse, error)
	GetBillingById(BillingId string) (models.Billing, error)
	UpdateBillingById(BillingId string, userId string, req models.UpdateBillingRequest) (models.Billing, error)
}

// BillingServiceImpl implements BillingService using GORM.
type BillingServiceImpl struct {
	db database.DatabaseManager
}

// NewBillingService creates a new instance of BillingService.
func NewBillingService(db *gorm.DB) BillingService {
	return &BillingServiceImpl{db: inst.InitDB(db)}
}

// CreateBilling creates a new billing record.
func (s *BillingServiceImpl) CreateBilling(req models.CreateBillingRequest, userId string) (models.BillingResponse, error) {
	var (
		user        models.User
		billingResp models.BillingResponse
	)

	Billing := models.Billing{
		ID:    utility.GenerateUUID(),
		Name:  req.Name,
		Price: req.Price,
	}

	if err := Billing.Create(s.db); err != nil {
		return billingResp, err
	}

	user, err := user.GetUserByID(s.db, userId)
	if err != nil {
		return billingResp, err
	}

	response := models.BillingResponse{
		BillingID: Billing.ID,
		Name:      Billing.Name,
		Price:     Billing.Price,
		CreatedAt: Billing.CreatedAt,
		UpdatedAt: Billing.UpdatedAt,
	}

	return response, nil
}

// DeleteBilling deletes a billing record by ID.
func (s *BillingServiceImpl) DeleteBilling(BillingId string, userId string) error {
	var Billing models.Billing

	Billing, err := Billing.CheckBillingExists(BillingId, s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("billing not found")
		}
		return err
	}

	return Billing.Delete(s.db)
}

// GetBillings retrieves paginated billing records.
func (s *BillingServiceImpl) GetBillings(c *gin.Context) (int, database.PaginationResponse, error) {
	var Billing models.Billing
	Billings, paginationResponse, err := Billing.GetAllBillings(s.db, c)
	if err != nil {
		return 0, paginationResponse, err
	}

	totalBillings := len(Billings)
	return totalBillings, paginationResponse, nil
}

// GetBillingById fetches a billing record by ID.
func (s *BillingServiceImpl) GetBillingById(BillingId string) (models.Billing, error) {
	var resp models.Billing
	resp, err := resp.GetBillingById(s.db, BillingId)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// UpdateBillingById updates a billing record by ID.
func (s *BillingServiceImpl) UpdateBillingById(BillingId string, userId string, req models.UpdateBillingRequest) (models.Billing, error) {
	var resp models.Billing

	resp, err := resp.CheckBillingExists(BillingId, s.db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, errors.New("billing not found")
		}
		return resp, err
	}

	_, err = resp.UpdateBillingById(s.db, req, BillingId)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
