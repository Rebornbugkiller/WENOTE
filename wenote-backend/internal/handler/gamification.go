package handler

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// GamificationHandler 游戏化处理器
type GamificationHandler struct {
	service *service.GamificationService
}

// NewGamificationHandler 创建游戏化处理器实例
func NewGamificationHandler() *GamificationHandler {
	return &GamificationHandler{
		service: service.NewGamificationService(),
	}
}

// GetStatus 获取游戏化状态
// GET /api/v1/gamification/status
func (h *GamificationHandler) GetStatus(c *gin.Context) {
	userID := c.GetUint64("userID")

	status, err := h.service.GetStatus(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, status)
}

// GetAchievements 获取所有成就及解锁状态
// GET /api/v1/gamification/achievements
func (h *GamificationHandler) GetAchievements(c *gin.Context) {
	userID := c.GetUint64("userID")

	achievements, err := h.service.GetAchievements(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": achievements})
}

// UpdateGoal 更新每日字符目标
// POST /api/v1/gamification/goal
func (h *GamificationHandler) UpdateGoal(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.UpdateGoalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := h.service.UpdateDailyGoal(userID, req.DailyCharGoal); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetReport 获取写作报告
// GET /api/v1/gamification/report?period=week
func (h *GamificationHandler) GetReport(c *gin.Context) {
	userID := c.GetUint64("userID")

	period := c.DefaultQuery("period", "week")
	if period != "week" && period != "month" {
		period = "week"
	}

	report, err := h.service.GetReport(userID, period)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, report)
}

// MarkAchievementNotified 标记成就已通知
// POST /api/v1/gamification/achievements/:id/notify
func (h *GamificationHandler) MarkAchievementNotified(c *gin.Context) {
	userID := c.GetUint64("userID")
	achievementID := c.Param("id")

	if achievementID == "" {
		response.BadRequest(c, "成就ID不能为空")
		return
	}

	if err := h.service.MarkAchievementNotified(userID, achievementID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
