package v1

import (
	"cryptoapi/internal/domain/services/cryptoanalysis/frequency"
	"cryptoapi/internal/domain/services/cryptoanalysis/polyalphabetic"
	"cryptoapi/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	freqService *frequency.Service
	polyService *polyalphabetic.Service
}

func NewController() *Controller {
	return &Controller{
		freqService: frequency.NewService(),
		polyService: polyalphabetic.NewService(),
	}
}

func (h *Controller) Init(r *gin.RouterGroup) {
	r.POST("/freqanalyze", h.AnalyseHandler)
	r.POST("/polyanalyze", h.PolyAnalyseHandler)
}

func (h *Controller) AnalyseHandler(c *gin.Context) {
	var req response.Response
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Lang != "en" && req.Lang != "ru" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported language, use 'en' or 'ru'"})
		return
	}
	result, err := h.freqService.Analyze(req.Text, req.Lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Controller) PolyAnalyseHandler(c *gin.Context) {
	var req response.Response
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Lang != "en" && req.Lang != "ru" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported language"})
		return
	}
	result := h.polyService.Analyze(req.Text, req.Lang)
	c.JSON(http.StatusOK, result)
}
