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

func CreateBilling(req models.CreateBillingRequest, db *gorm.DB, userId string) (models.BillingResponse, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)
	var (
		user        models.User
		billingResp models.BillingResponse
	)
	Billing := models.Billing{
		ID:    utility.GenerateUUID(),
		Name:  req.Name,
		Price: req.Price,
	}

	err := Billing.Create(pdb)

	if err != nil {
		return billingResp, err
	}

	user, err = user.GetUserByID(pdb, userId)

	if err != nil {
		return billingResp, err
	}

	response := models.BillingResponse{
		BillingID: Billing.ID,
		Name:      Billing.Name,
		Price:     Billing.Price,
		CreatedAt: Billing.CreatedAt,
		UpdatedAt: billingResp.UpdatedAt,
	}

	return response, nil
}

func DeleteBilling(BillingId string, userId string, db *gorm.DB) error {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	var Billing models.Billing

	Billing, err := Billing.CheckBillingExists(BillingId, pdb)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Billing not found")
		}
		return err
	}

	return Billing.Delete(pdb)
}

func GetBillings(db *gorm.DB, c *gin.Context) (int, database.PaginationResponse, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	var (
		Billing models.Billing
	)
	Billings, paginationResponse, err := Billing.GetAllBillings(pdb, c)

	if err != nil {
		return 0, paginationResponse, err
	}

	total_billings := len(Billings)

	return total_billings, paginationResponse, nil
}

func GetBillingById(BillingId string, db *gorm.DB) (models.Billing, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	var (
		resp models.Billing
	)

	resp, err := resp.GetBillingById(pdb, BillingId)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func UpdateBillingById(BillingId string, userId string, req models.UpdateBillingRequest, db *gorm.DB) (models.Billing, error) {
	// instance of Postgresql db
	pdb := inst.InitDB(db)

	var (
		resp models.Billing
	)

	resp, err := resp.CheckBillingExists(BillingId, pdb)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, errors.New("billing not found")
		}
		return resp, err
	}

	_, err = resp.UpdateBillingById(pdb, req, BillingId)

	if err != nil {
		return resp, err
	}

	return resp, nil
}
