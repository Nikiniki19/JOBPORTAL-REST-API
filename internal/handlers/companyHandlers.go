package handlers

import (
	"encoding/json"
	"job-portal-api/internal/middlewares"
	"job-portal-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

// Create company API
func (h *handler) createCom(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg1": http.StatusText(http.StatusInternalServerError)})
		return
	}
	var newComp models.Company
	err := json.NewDecoder(c.Request.Body).Decode(&newComp)
	if err != nil {
		log.Error().Err(err).Msg("error in decoding")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg2": http.StatusText(http.StatusInternalServerError)})
		return
	}
	validate := validator.New()
	err = validate.Struct(&newComp)
	if err != nil {
		log.Error().Err(err).Msg("validation error")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg3":http.StatusText(http.StatusBadRequest)})
		return
	}
	comp, err := h.s.AddCompanyDetails(ctx, newComp)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("user login problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg4": "user not found"})
		return
	}
	c.JSON(http.StatusOK, comp)

}

// Getting all the companies API
func (h *handler) getAllTheCompanies(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	comp, err := h.s.ViewAllCompanies(ctx)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return

	}
	c.JSON(http.StatusOK, comp)
}

// Viewing a company by fetching id API
func (h *handler) viewCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Str("traceId", traceId).Msg("trace id missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	id, error := strconv.ParseUint(c.Param("cid"), 10, 64)
	if error != nil {
		log.Error().Str("traceId", traceId).Msg("company id invalid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": http.StatusText(http.StatusBadRequest)})
		return
	}
	comp, err := h.s.ViewCompanyDetails(ctx, id)
	if err != nil {
		log.Error().Str("traceId", traceId).Msg("companies not founds")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusOK, comp)

}
