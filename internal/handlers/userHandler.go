package handlers

import (
	"encoding/json"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/middlewares"
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// Handler Struct
type handler struct {
	s services.UserService
	a auth.Authentication
}

// Registration API
func (h *handler) Registration(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	var nu models.NewUser
	err := json.NewDecoder(c.Request.Body).Decode(&nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	validate := validator.New()
	err = validate.Struct(nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": http.StatusText(http.StatusBadRequest)})
		return
	}
	usr, err := h.s.Signup(ctx, nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("user signup problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusOK, usr)

}

// Login API
func (h *handler) Signin(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var login models.Login

	err := json.NewDecoder(c.Request.Body).Decode(&login)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// Create a new validator and validate the login variable
	validate := validator.New()
	err = validate.Struct(login)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Email and Password"})
		return
	}

	// Attempt to authenticate the user with the email and password
	claims, err := h.s.Login(ctx, login.Email, login.Password)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login failed"})
		return
	}

	// Generate a new token and put it in the Token field of the token struct
	tkn, err := h.a.GenerateToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("generating token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// If everything goes right, respond with the token
	c.JSON(http.StatusOK, gin.H{"token":tkn})

}
func(h *handler)ForgotPassword(c *gin.Context){
	ctx:=c.Request.Context()
	traceid,ok:=ctx.Value(middlewares.TraceIdKey).(string)
	if !ok{
		log.Error().Msg("traceid missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"error":http.StatusText(http.StatusInternalServerError)})
		return
	}
	 var fp models.ForgotPassword
	 err:=json.NewDecoder(c.Request.Body).Decode(&fp)
	 if err!=nil{
		log.Error().Str("traceid",traceid).Msg("error in decoding")
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"error":http.StatusText(http.StatusBadRequest)})
		return
	 }

	 validate:=validator.New()
	 err=validate.Struct(fp)
	 if err!=nil{
		log.Error().Str("traceid",traceid).Msg("error in validating")
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"error":http.StatusText(http.StatusBadRequest)})
		return
	 }
      
	valid,err:=h.s.OTPGeneration()
	if err!=nil{
	   log.Error().Str("traceid",traceid).Msg("error in generating otp")
	   c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"error":http.StatusText(http.StatusInternalServerError)})
	   return
	}





}
