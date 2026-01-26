package repo

import (
	"time"
	"wenote-backend/internal/model"
)

// GamificationRepo 游戏化数据访问层
type GamificationRepo struct{}

// NewGamificationRepo 创建游戏化仓库实例
func NewGamificationRepo() *GamificationRepo {
	return &GamificationRepo{}
}

// GetOrCreateUserGamification 获取或创建用户游戏化数据
func (r *GamificationRepo) GetOrCreateUserGamification(userID uint64) (*model.UserGamification, error) {
	var ug model.UserGamification
	result := DB.Where("user_id = ?", userID).First(&ug)
	if result.Error != nil {
		// 不存在则创建
		ug = model.UserGamification{
			UserID:        userID,
			DailyCharGoal: 500,
		}
		if err := DB.Create(&ug).Error; err != nil {
			return nil, err
		}
	}
	return &ug, nil
}

// UpdateUserGamification 更新用户游戏化数据
func (r *GamificationRepo) UpdateUserGamification(ug *model.UserGamification) error {
	return DB.Save(ug).Error
}

// GetAllAchievements 获取所有成就定义
func (r *GamificationRepo) GetAllAchievements() ([]model.Achievement, error) {
	var achievements []model.Achievement
	err := DB.Find(&achievements).Error
	return achievements, err
}

// GetUserAchievements 获取用户已解锁的成就
func (r *GamificationRepo) GetUserAchievements(userID uint64) ([]model.UserAchievement, error) {
	var userAchievements []model.UserAchievement
	err := DB.Where("user_id = ?", userID).Find(&userAchievements).Error
	return userAchievements, err
}

// GetUnnotifiedAchievements 获取用户未通知的成就
func (r *GamificationRepo) GetUnnotifiedAchievements(userID uint64) ([]model.Achievement, error) {
	var achievements []model.Achievement
	err := DB.Raw(`
		SELECT a.* FROM achievements a
		INNER JOIN user_achievements ua ON a.id = ua.achievement_id
		WHERE ua.user_id = ? AND ua.notified = false
	`, userID).Scan(&achievements).Error
	return achievements, err
}

// UnlockAchievement 解锁成就
func (r *GamificationRepo) UnlockAchievement(userID uint64, achievementID string) error {
	ua := model.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
	}
	// 使用 FirstOrCreate 避免重复解锁
	return DB.Where("user_id = ? AND achievement_id = ?", userID, achievementID).FirstOrCreate(&ua).Error
}

// MarkAchievementNotified 标记成就已通知
func (r *GamificationRepo) MarkAchievementNotified(userID uint64, achievementID string) error {
	return DB.Model(&model.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Update("notified", true).Error
}

// HasAchievement 检查用户是否已解锁某成就
func (r *GamificationRepo) HasAchievement(userID uint64, achievementID string) (bool, error) {
	var count int64
	err := DB.Model(&model.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error
	return count > 0, err
}

// GetTotalNotes 获取用户总笔记数
func (r *GamificationRepo) GetTotalNotes(userID uint64) (int64, error) {
	var count int64
	err := DB.Model(&model.Note{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&count).Error
	return count, err
}

// GetTotalChars 获取用户总字符数
func (r *GamificationRepo) GetTotalChars(userID uint64) (int64, error) {
	var result struct {
		Total int64
	}
	err := DB.Raw("SELECT COALESCE(SUM(LENGTH(content)), 0) as total FROM notes WHERE user_id = ? AND deleted_at IS NULL", userID).Scan(&result).Error
	return result.Total, err
}

// GetActiveDaysInPeriod 获取指定时间段内的活跃天数
func (r *GamificationRepo) GetActiveDaysInPeriod(userID uint64, startDate, endDate time.Time) (int, error) {
	var count int64
	err := DB.Raw(`
		SELECT COUNT(DISTINCT DATE(created_at)) as count
		FROM notes
		WHERE user_id = ? AND deleted_at IS NULL
		AND created_at >= ? AND created_at < ?
	`, userID, startDate, endDate).Scan(&count).Error
	return int(count), err
}

// GetNotesCreatedInPeriod 获取指定时间段内创建的笔记数
func (r *GamificationRepo) GetNotesCreatedInPeriod(userID uint64, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := DB.Model(&model.Note{}).
		Where("user_id = ? AND deleted_at IS NULL AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetCharsWrittenInPeriod 获取指定时间段内写的字符数
func (r *GamificationRepo) GetCharsWrittenInPeriod(userID uint64, startDate, endDate time.Time) (int64, error) {
	var result struct {
		Total int64
	}
	err := DB.Raw(`
		SELECT COALESCE(SUM(LENGTH(content)), 0) as total
		FROM notes
		WHERE user_id = ? AND deleted_at IS NULL
		AND created_at >= ? AND created_at < ?
	`, userID, startDate, endDate).Scan(&result).Error
	return result.Total, err
}

// GetDailyStatsInPeriod 获取指定时间段内的每日统计
func (r *GamificationRepo) GetDailyStatsInPeriod(userID uint64, startDate, endDate time.Time) ([]model.DailyStat, error) {
	var results []struct {
		Date  string
		Notes int64
		Chars int64
	}

	err := DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as notes, COALESCE(SUM(LENGTH(content)), 0) as chars
		FROM notes
		WHERE user_id = ? AND deleted_at IS NULL
		AND created_at >= ? AND created_at < ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, userID, startDate, endDate).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为 DailyStat
	stats := make([]model.DailyStat, len(results))
	for i, r := range results {
		stats[i] = model.DailyStat{
			Date:  r.Date,
			Notes: r.Notes,
			Chars: r.Chars,
		}
	}

	return stats, nil
}

// GetAchievementsUnlockedInPeriod 获取指定时间段内解锁的成就
func (r *GamificationRepo) GetAchievementsUnlockedInPeriod(userID uint64, startDate, endDate time.Time) ([]model.Achievement, error) {
	var achievements []model.Achievement
	err := DB.Raw(`
		SELECT a.* FROM achievements a
		INNER JOIN user_achievements ua ON a.id = ua.achievement_id
		WHERE ua.user_id = ? AND ua.unlocked_at >= ? AND ua.unlocked_at < ?
	`, userID, startDate, endDate).Scan(&achievements).Error
	return achievements, err
}
