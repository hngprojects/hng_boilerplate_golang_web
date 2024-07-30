package test_blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/auth"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/blog"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	tst "github.com/hngprojects/hng_boilerplate_golang_web/tests"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func initialise(currUUID string, t *testing.T, r *gin.Engine, db *storage.Database, user auth.Controller, blog blog.Controller, status bool) (string, string) {
	userSignUpData := models.CreateUserRequestModel{
		Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
		PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		UserName:    fmt.Sprintf("test_username%v", currUUID),
	}
	loginData := models.LoginRequestModel{
		Email:    userSignUpData.Email,
		Password: userSignUpData.Password,
	}

	tst.SignupUser(t, r, user, userSignUpData, status)

	token := tst.GetLoginToken(t, r, user, loginData)

	blogCreationData := models.CreateBlogRequest{
		Title:   fmt.Sprintf("Blog %s", currUUID),
		Content: fmt.Sprintf("testuser%s", currUUID),
	}

	blogID := CreateBlog(t, r, db, blog, blogCreationData, token)

	return blogID, token
}

func CreateBlog(t *testing.T, r *gin.Engine, db *storage.Database, blog blog.Controller, blogData models.CreateBlogRequest, token string) string {
	var (
		blogPath = "/api/v1/blogs"
		blogURI  = url.URL{Path: blogPath}
	)
	blogUrl := r.Group(fmt.Sprintf("%v", "/api/v1"), middleware.Authorize(db.Postgresql, models.RoleIdentity.SuperAdmin, models.RoleIdentity.User))
	{
		blogUrl.POST("/blogs", blog.CreateBlog)
	}
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(blogData)
	req, err := http.NewRequest(http.MethodPost, blogURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	//get the response
	data := tst.ParseResponse(rr)
	dataM := data["data"].(map[string]interface{})
	blogID := dataM["id"].(string)
	return blogID
}
