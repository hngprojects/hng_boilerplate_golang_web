package python

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) RunPythonTestsHandler(c *gin.Context) {
	cmd := exec.Command("python3", "scripts/compare_test.py")
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if err != nil {

		base.Logger.Error("Failed to run Python script: %s", err)
		base.Logger.Error("Error output: %s", errb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run tests", "details": errb.String()})
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", outb.Bytes())
}
