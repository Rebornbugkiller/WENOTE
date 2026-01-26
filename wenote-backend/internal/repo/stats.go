package repo

import (
	"time"
	"wenote-backend/internal/model"
)

// StatsRepo 统计数据访问层
type StatsRepo struct{}

// NewStatsRepo 创建统计仓库实例
func NewStatsRepo() *StatsRepo {
	return &StatsRepo{}
}

// GetOverview 获取统计概览
func (r *StatsRepo) GetOverview(userID uint64) (*model.StatsOverview, error) {
	var stats model.StatsOverview

	// 获取本周开始时间
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	// 总笔记数
	DB.Model(&model.Note{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&stats.TotalNotes)

	// 总笔记本数
	DB.Model(&model.Notebook{}).Where("user_id = ?", userID).Count(&stats.TotalNotebooks)

	// 总标签数
	DB.Model(&model.Tag{}).Where("user_id = ?", userID).Count(&stats.TotalTags)

	// 本周新增笔记数
	DB.Model(&model.Note{}).Where("user_id = ? AND deleted_at IS NULL AND created_at >= ?", userID, weekStart).Count(&stats.ThisWeekNotes)

	// 总字数（统计content字段长度）
	var totalWords struct {
		Total int64
	}
	DB.Raw("SELECT COALESCE(SUM(LENGTH(content)), 0) as total FROM notes WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&totalWords)
	stats.TotalWords = totalWords.Total

	// 本周新增字数
	var weekWords struct {
		Total int64
	}
	DB.Raw("SELECT COALESCE(SUM(LENGTH(content)), 0) as total FROM notes WHERE user_id = ? AND deleted_at IS NULL AND created_at >= ?", userID, weekStart).Scan(&weekWords)
	stats.ThisWeekWords = weekWords.Total

	return &stats, nil
}

// GetTrendData 获取趋势数据
func (r *StatsRepo) GetTrendData(userID uint64, days int) ([]model.TrendData, error) {
	var results []model.TrendData

	// 计算开始日期
	startDate := time.Now().AddDate(0, 0, -days+1)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// 查询每天的笔记数量
	err := DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM notes
		WHERE user_id = ? AND deleted_at IS NULL AND created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, userID, startDate).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 填充缺失的日期（确保每一天都有数据点）
	dateMap := make(map[string]int64)
	for _, r := range results {
		dateMap[r.Date] = r.Count
	}

	fullResults := make([]model.TrendData, 0, days)
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		count := dateMap[dateStr]
		fullResults = append(fullResults, model.TrendData{
			Date:  dateStr,
			Count: count,
		})
	}

	return fullResults, nil
}

// GetTagStats 获取标签统计（TOP N）
func (r *StatsRepo) GetTagStats(userID uint64, limit int) ([]model.TagStat, error) {
	var results []model.TagStat

	err := DB.Raw(`
		SELECT t.name as tag_name, t.color, COUNT(nt.note_id) as count
		FROM tags t
		LEFT JOIN note_tags nt ON t.id = nt.tag_id
		LEFT JOIN notes n ON nt.note_id = n.id AND n.deleted_at IS NULL
		WHERE t.user_id = ?
		GROUP BY t.id, t.name, t.color
		HAVING count > 0
		ORDER BY count DESC
		LIMIT ?
	`, userID, limit).Scan(&results).Error

	return results, err
}

// GetNotebookStats 获取笔记本统计
func (r *StatsRepo) GetNotebookStats(userID uint64) ([]model.NotebookStat, error) {
	var results []model.NotebookStat

	err := DB.Raw(`
		SELECT nb.name as notebook_name, COUNT(n.id) as count
		FROM notebooks nb
		LEFT JOIN notes n ON nb.id = n.notebook_id AND n.deleted_at IS NULL
		WHERE nb.user_id = ?
		GROUP BY nb.id, nb.name
		ORDER BY count DESC
	`, userID).Scan(&results).Error

	return results, err
}



