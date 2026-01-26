package handler

import (
	"strconv"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// StatsHandler 统计处理器
type StatsHandler struct {
	service *service.StatsService
}

// NewStatsHandler 创建统计处理器实例
func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		service: service.NewStatsService(),
	}
}

// GetOverview 获取统计概览
// GET /api/v1/stats/overview
func (h *StatsHandler) GetOverview(c *gin.Context) {
	userID := c.GetUint64("userID")

	overview, err := h.service.GetOverview(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, overview)
}

// GetTrendData 获取趋势数据
// GET /api/v1/stats/trend?days=7
func (h *StatsHandler) GetTrendData(c *gin.Context) {
	userID := c.GetUint64("userID")

	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 7
	}

	trendData, err := h.service.GetTrendData(userID, days)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": trendData})
}

// GetTagStats 获取标签统计
// GET /api/v1/stats/tags?limit=10
func (h *StatsHandler) GetTagStats(c *gin.Context) {
	userID := c.GetUint64("userID")

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	tagStats, err := h.service.GetTagStats(userID, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": tagStats})
}

// GetNotebookStats 获取笔记本统计
// GET /api/v1/stats/notebooks
func (h *StatsHandler) GetNotebookStats(c *gin.Context) {
	userID := c.GetUint64("userID")

	notebookStats, err := h.service.GetNotebookStats(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": notebookStats})
}



