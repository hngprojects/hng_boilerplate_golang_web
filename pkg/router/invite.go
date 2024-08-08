package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/controller/invite"

	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/middleware"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func Invite(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	invite := invite.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	inviteUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db.Postgresql))
	{
		{
			inviteUrl.POST("/invite/create", invite.CreateInvite)
			inviteUrl.POST("/organization/send-invite", middleware.RateLimiter(), invite.PostInvite)
			inviteUrl.POST("/invite/accept", invite.PostAcceptInvite)
		}


		{
			inviteUrl.GET("/organization/invitations", invite.GetInvites)
			inviteUrl.GET("/invite/accept/:t", invite.GetAcceptInvite)
		}

	}
	return r
}