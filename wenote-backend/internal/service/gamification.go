package service

import (
	"time"
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
)

// GamificationService 游戏化服务
type GamificationService struct {
	repo *repo.GamificationRepo
}

// NewGamificationService 创建游戏化服务实例
func NewGamificationService() *GamificationService {
	return &GamificationService{
		repo: repo.NewGamificationRepo(),
	}
}

// GetStatus 获取用户游戏化状态
func (s *GamificationService) GetStatus(userID uint64) (*model.GamificationStatus, error) {
	// 获取用户游戏化数据
	ug, err := s.repo.GetOrCreateUserGamification(userID)
	if err != nil {
		return nil, err
	}

	// 检查并更新日期相关数据
	today := time.Now().Truncate(24 * time.Hour)
	s.checkAndResetDaily(ug, today)

	// 获取统计数据
	totalNotes, _ := s.repo.GetTotalNotes(userID)
	totalChars, _ := s.repo.GetTotalChars(userID)

	// 获取未通知的成就
	newAchievements, _ := s.repo.GetUnnotifiedAchievements(userID)
	achievementsWithStatus := make([]model.AchievementWithStatus, len(newAchievements))
	for i, a := range newAchievements {
		achievementsWithStatus[i] = model.AchievementWithStatus{
			Achievement: a,
			Unlocked:    true,
		}
	}

	// 计算目标进度
	goalProgress := float64(0)
	if ug.DailyCharGoal > 0 {
		goalProgress = float64(ug.TodayChars) / float64(ug.DailyCharGoal)
	}

	// 判断连续天数是否有风险（今天还没有活动）
	streakAtRisk := false
	if ug.CurrentStreak > 0 && (ug.LastActiveDate == nil || !isSameDay(*ug.LastActiveDate, today)) {
		streakAtRisk = true
	}

	return &model.GamificationStatus{
		CurrentStreak:   ug.CurrentStreak,
		LongestStreak:   ug.LongestStreak,
		LastActiveDate:  ug.LastActiveDate,
		StreakAtRisk:    streakAtRisk,
		DailyCharGoal:   ug.DailyCharGoal,
		TodayChars:      ug.TodayChars,
		GoalProgress:    goalProgress,
		GoalCompleted:   ug.TodayChars >= ug.DailyCharGoal,
		TotalNotes:      totalNotes,
		TotalChars:      totalChars,
		NewAchievements: achievementsWithStatus,
	}, nil
}

// UpdateActivity 更新用户活动（创建/编辑笔记时调用）
func (s *GamificationService) UpdateActivity(userID uint64, charsDelta int64) error {
	ug, err := s.repo.GetOrCreateUserGamification(userID)
	if err != nil {
		return err
	}

	today := time.Now().Truncate(24 * time.Hour)

	// 检查并重置每日数据
	s.checkAndResetDaily(ug, today)

	// 更新今日字符数
	if charsDelta > 0 {
		ug.TodayChars += int(charsDelta)
	}

	// 检查是否完成每日目标
	goalWasNotCompleted := ug.TodayChars-int(charsDelta) < ug.DailyCharGoal
	goalNowCompleted := ug.TodayChars >= ug.DailyCharGoal
	if goalWasNotCompleted && goalNowCompleted {
		ug.GoalsCompleted++
	}

	// 更新连续天数
	if ug.LastActiveDate == nil || !isSameDay(*ug.LastActiveDate, today) {
		yesterday := today.AddDate(0, 0, -1)
		if ug.LastActiveDate != nil && isSameDay(*ug.LastActiveDate, yesterday) {
			// 昨天有活动，连续天数+1
			ug.CurrentStreak++
		} else if ug.LastActiveDate == nil || !isSameDay(*ug.LastActiveDate, today) {
			// 不是连续的，重置为1
			ug.CurrentStreak = 1
		}

		// 更新最长连续天数
		if ug.CurrentStreak > ug.LongestStreak {
			ug.LongestStreak = ug.CurrentStreak
		}

		// 更新最后活跃日期
		ug.LastActiveDate = &today
	}

	// 保存更新
	if err := s.repo.UpdateUserGamification(ug); err != nil {
		return err
	}

	// 检查并解锁成就
	s.checkAndUnlockAchievements(userID, ug)

	return nil
}

// checkAndResetDaily 检查并重置每日数据
func (s *GamificationService) checkAndResetDaily(ug *model.UserGamification, today time.Time) {
	if ug.TodayDate == nil || !isSameDay(*ug.TodayDate, today) {
		ug.TodayChars = 0
		ug.TodayDate = &today
	}
}

// checkAndUnlockAchievements 检查并解锁成就
func (s *GamificationService) checkAndUnlockAchievements(userID uint64, ug *model.UserGamification) {
	totalNotes, _ := s.repo.GetTotalNotes(userID)
	totalChars, _ := s.repo.GetTotalChars(userID)

	// 笔记数量成就
	noteAchievements := []struct {
		id        string
		threshold int64
	}{
		{"first_note", 1},
		{"notes_10", 10},
		{"notes_50", 50},
		{"notes_100", 100},
		{"notes_500", 500},
	}
	for _, a := range noteAchievements {
		if totalNotes >= a.threshold {
			s.tryUnlockAchievement(userID, a.id)
		}
	}

	// 连续天数成就
	streakAchievements := []struct {
		id        string
		threshold int
	}{
		{"streak_3", 3},
		{"streak_7", 7},
		{"streak_30", 30},
		{"streak_100", 100},
	}
	for _, a := range streakAchievements {
		if ug.CurrentStreak >= a.threshold {
			s.tryUnlockAchievement(userID, a.id)
		}
	}

	// 字符数成就
	charAchievements := []struct {
		id        string
		threshold int64
	}{
		{"chars_1k", 1000},
		{"chars_10k", 10000},
		{"chars_50k", 50000},
		{"chars_100k", 100000},
	}
	for _, a := range charAchievements {
		if totalChars >= a.threshold {
			s.tryUnlockAchievement(userID, a.id)
		}
	}

	// 目标完成成就
	goalAchievements := []struct {
		id        string
		threshold int
	}{
		{"goal_first", 1},
		{"goal_7", 7},
		{"goal_30", 30},
	}
	for _, a := range goalAchievements {
		if ug.GoalsCompleted >= a.threshold {
			s.tryUnlockAchievement(userID, a.id)
		}
	}
}

// tryUnlockAchievement 尝试解锁成就（如果未解锁）
func (s *GamificationService) tryUnlockAchievement(userID uint64, achievementID string) {
	has, _ := s.repo.HasAchievement(userID, achievementID)
	if !has {
		s.repo.UnlockAchievement(userID, achievementID)
	}
}

// GetAchievements 获取所有成就及用户解锁状态
func (s *GamificationService) GetAchievements(userID uint64) ([]model.AchievementWithStatus, error) {
	// 获取所有成就
	achievements, err := s.repo.GetAllAchievements()
	if err != nil {
		return nil, err
	}

	// 获取用户已解锁的成就
	userAchievements, err := s.repo.GetUserAchievements(userID)
	if err != nil {
		return nil, err
	}

	// 构建已解锁成就的 map
	unlockedMap := make(map[string]time.Time)
	for _, ua := range userAchievements {
		unlockedMap[ua.AchievementID] = ua.UnlockedAt
	}

	// 组合结果
	result := make([]model.AchievementWithStatus, len(achievements))
	for i, a := range achievements {
		result[i] = model.AchievementWithStatus{
			Achievement: a,
			Unlocked:    false,
		}
		if unlockedAt, ok := unlockedMap[a.ID]; ok {
			result[i].Unlocked = true
			result[i].UnlockedAt = &unlockedAt
		}
	}

	return result, nil
}

// UpdateDailyGoal 更新每日字符目标
func (s *GamificationService) UpdateDailyGoal(userID uint64, goal int) error {
	ug, err := s.repo.GetOrCreateUserGamification(userID)
	if err != nil {
		return err
	}

	ug.DailyCharGoal = goal
	return s.repo.UpdateUserGamification(ug)
}

// MarkAchievementNotified 标记成就已通知
func (s *GamificationService) MarkAchievementNotified(userID uint64, achievementID string) error {
	return s.repo.MarkAchievementNotified(userID, achievementID)
}

// GetReport 获取写作报告
func (s *GamificationService) GetReport(userID uint64, period string) (*model.WritingReport, error) {
	now := time.Now()
	var startDate, endDate, prevStartDate, prevEndDate time.Time

	if period == "month" {
		// 本月
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
		// 上月
		prevStartDate = startDate.AddDate(0, -1, 0)
		prevEndDate = startDate
	} else {
		// 本周（默认）
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startDate = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 0, 7)
		// 上周
		prevStartDate = startDate.AddDate(0, 0, -7)
		prevEndDate = startDate
	}

	// 获取本期数据
	notesCreated, _ := s.repo.GetNotesCreatedInPeriod(userID, startDate, endDate)
	charsWritten, _ := s.repo.GetCharsWrittenInPeriod(userID, startDate, endDate)
	activeDays, _ := s.repo.GetActiveDaysInPeriod(userID, startDate, endDate)
	dailyStats, _ := s.repo.GetDailyStatsInPeriod(userID, startDate, endDate)
	achievementsEarned, _ := s.repo.GetAchievementsUnlockedInPeriod(userID, startDate, endDate)

	// 获取上期数据（用于对比）
	prevNotes, _ := s.repo.GetNotesCreatedInPeriod(userID, prevStartDate, prevEndDate)
	prevChars, _ := s.repo.GetCharsWrittenInPeriod(userID, prevStartDate, prevEndDate)
	prevActiveDays, _ := s.repo.GetActiveDaysInPeriod(userID, prevStartDate, prevEndDate)

	// 获取用户的每日目标，用于判断每天是否达标
	ug, _ := s.repo.GetOrCreateUserGamification(userID)
	goalsMet := 0
	for i := range dailyStats {
		if dailyStats[i].Chars >= int64(ug.DailyCharGoal) {
			dailyStats[i].GoalMet = true
			goalsMet++
		}
	}

	return &model.WritingReport{
		Period:             period,
		StartDate:          startDate.Format("2006-01-02"),
		EndDate:            endDate.AddDate(0, 0, -1).Format("2006-01-02"),
		NotesCreated:       notesCreated,
		CharsWritten:       charsWritten,
		ActiveDays:         activeDays,
		GoalsMet:           goalsMet,
		NotesDelta:         notesCreated - prevNotes,
		CharsDelta:         charsWritten - prevChars,
		ActiveDaysDelta:    activeDays - prevActiveDays,
		AchievementsEarned: achievementsEarned,
		DailyStats:         dailyStats,
	}, nil
}

// isSameDay 判断两个时间是否是同一天
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
