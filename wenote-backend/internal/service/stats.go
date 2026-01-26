package service

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
)

// StatsService 统计服务
type StatsService struct {
	repo *repo.StatsRepo
}

// NewStatsService 创建统计服务实例
func NewStatsService() *StatsService {
	return &StatsService{
		repo: repo.NewStatsRepo(),
	}
}

// GetOverview 获取统计概览
func (s *StatsService) GetOverview(userID uint64) (*model.StatsOverview, error) {
	return s.repo.GetOverview(userID)
}

// GetTrendData 获取趋势数据
func (s *StatsService) GetTrendData(userID uint64, days int) ([]model.TrendData, error) {
	// 限制天数范围
	if days <= 0 {
		days = 7
	}
	if days > 90 {
		days = 90
	}
	return s.repo.GetTrendData(userID, days)
}

// GetTagStats 获取标签统计
func (s *StatsService) GetTagStats(userID uint64, limit int) ([]model.TagStat, error) {
	// 限制返回数量
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}
	return s.repo.GetTagStats(userID, limit)
}

// GetNotebookStats 获取笔记本统计
func (s *StatsService) GetNotebookStats(userID uint64) ([]model.NotebookStat, error) {
	return s.repo.GetNotebookStats(userID)
}



