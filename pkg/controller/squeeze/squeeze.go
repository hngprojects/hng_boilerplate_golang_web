package squeeze

import (
	

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	DB        *storage.Database
	Logger    *utility.Logger
	Validator *validator.Validate
}

func (base *Controller) Create(c *gin.Context) {

}
