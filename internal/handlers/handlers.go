package handlers

import (
	"fmt"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/middlewares"
	"job-portal-api/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// API func
func API(a auth.Authentication, s services.UserService) *gin.Engine {

	r := gin.New()

	// Attempt to create new middleware with authentication
	// Here, *auth.Auth passed as a parameter will be used to set up the middleware
	m, _ := middlewares.NewMid(a)
	h := handler{
		a: a,
		s: s,
	}

	r.Use(m.LoggerMiddleware(), gin.Recovery())

	//Endpoints call
	r.GET("/check", m.AuthenticationMiddleware(check))
	//users endpoint
	r.POST("/signup", h.Registration)
	r.POST("/login", h.Signin)
	//company endpoint
	r.POST("/createCompany", m.AuthenticationMiddleware(h.createCom))
	r.GET("/getallcompanies", m.AuthenticationMiddleware(h.getAllTheCompanies))
	r.GET("/getacompany/:cid", m.AuthenticationMiddleware(h.viewCompany))
	//jobs endpoint
	r.POST("/companies/:cid", m.AuthenticationMiddleware(h.postJob))
	r.GET("/jobs/:CompanyId", m.AuthenticationMiddleware(h.getJobsFromCompany))
	r.GET("/jobs", m.AuthenticationMiddleware(h.getAllJobs))
	r.GET("/jobs/jid", m.AuthenticationMiddleware(h.getOneJob))

	r.POST("/process/applications", m.AuthenticationMiddleware(h.processApplications))
	r.POST("/forget",h.ForgotPassword)
	r.POST("/password",h.SetNewPassword)


	return r
}

// Checking whether the user is  there or not
func check(c *gin.Context) {
	time.Sleep(time.Second * 3)
	select {
	case <-c.Request.Context().Done():
		fmt.Println("user not there")
		return
	default:
		c.JSON(http.StatusOK, gin.H{"msg": "statusOk"})

	}

}
