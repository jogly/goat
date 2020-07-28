package handler

import (
	"net/http"

	"github.com/banditml/goat/header"
	"github.com/banditml/goat/model"
	"github.com/banditml/goat/route"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func NewCampaignHandler(db *gorm.DB, logger *zap.Logger) route.HandlerInterface {
	return &CampaignHandler{
		l:  logger,
		db: db,
	}
}

type CampaignHandler struct {
	route.Base

	l  *zap.Logger
	db *gorm.DB
}

func (p *CampaignHandler) Resource() string {
	return "campaign"
}

func (p *CampaignHandler) Get(c *gin.Context) {
	account := c.GetHeader(header.BanditId)
	campaign := new(model.Campaign)
	if err := p.db.First(campaign, "account = ?", account).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		panic(err)
	}
	c.JSON(200, campaign)
}

func (p *CampaignHandler) Post(c *gin.Context) {
	account := c.GetHeader(header.BanditId)
	campaign := new(model.Campaign)
	if err := c.ShouldBindJSON(campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	campaign.Account = account
	err := p.db.Where(model.Campaign{Account: account}).
		Assign(*campaign).
		FirstOrCreate(campaign).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"campaign": campaign,
	})
}
