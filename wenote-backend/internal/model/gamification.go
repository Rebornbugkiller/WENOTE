package model

import (
	"time"
)

// ========== 游戏化模型 ==========

// UserGamification 用户游戏化数据
// 存储用户的连续天数、每日目标等游戏化状态
type UserGamification struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64     `gorm:"uniqueIndex;not null" json:"user_id"`
	CurrentStreak  int        `gorm:"default:0" json:"current_streak"`   // 当前连续天数
	LongestStreak  int        `gorm:"default:0" json:"longest_streak"`   // 最长连续天数
	LastActiveDate *time.Time `gorm:"type:date" json:"last_active_date"` // 最后活跃日期
	DailyCharGoal  int        `gorm:"default:500" json:"daily_char_goal"` // 每日字符目标
	TodayChars     int        `gorm:"default:0" json:"today_chars"`       // 今日字符数
	TodayDate      *time.Time `gorm:"type:date" json:"today_date"`        // today_chars 对应的日期
	GoalsCompleted int        `gorm:"default:0" json:"goals_completed"`   // 累计完成目标次数
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserGamification) TableName() string {
	return "user_gamification"
}

// Achievement 成就定义
// 预置的成就列表，存储在数据库中
type Achievement struct {
	ID            string `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Name          string `gorm:"type:varchar(100);not null" json:"name"`
	NameZh        string `gorm:"type:varchar(100);not null" json:"name_zh"`
	Description   string `gorm:"type:varchar(255)" json:"description"`
	DescriptionZh string `gorm:"type:varchar(255)" json:"description_zh"`
	Icon          string `gorm:"type:varchar(10)" json:"icon"`      // emoji
	Category      string `gorm:"type:varchar(50)" json:"category"`  // notes/streak/words/goals
	Threshold     int    `gorm:"default:0" json:"threshold"`        // 解锁阈值
	Rarity        string `gorm:"type:varchar(20);default:'common'" json:"rarity"` // common/rare/epic/legendary
}

func (Achievement) TableName() string {
	return "achievements"
}

// UserAchievement 用户已解锁的成就
type UserAchievement struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"index;not null" json:"user_id"`
	AchievementID string    `gorm:"type:varchar(50);not null" json:"achievement_id"`
	UnlockedAt    time.Time `gorm:"autoCreateTime" json:"unlocked_at"`
	Notified      bool      `gorm:"default:false" json:"notified"` // 是否已通知用户
}

func (UserAchievement) TableName() string {
	return "user_achievements"
}

// ========== 请求/响应 DTO ==========

// GamificationStatus 游戏化状态响应
// 用于 GET /api/v1/gamification/status
type GamificationStatus struct {
	// 连续天数
	CurrentStreak  int        `json:"current_streak"`
	LongestStreak  int        `json:"longest_streak"`
	LastActiveDate *time.Time `json:"last_active_date"`
	StreakAtRisk   bool       `json:"streak_at_risk"` // 今天还没有活动，连续天数有风险

	// 每日目标
	DailyCharGoal int     `json:"daily_char_goal"`
	TodayChars    int     `json:"today_chars"`
	GoalProgress  float64 `json:"goal_progress"` // 0.0 ~ 1.0+
	GoalCompleted bool    `json:"goal_completed"`

	// 统计
	TotalNotes int64 `json:"total_notes"`
	TotalChars int64 `json:"total_chars"`

	// 新解锁的成就（未通知）
	NewAchievements []AchievementWithStatus `json:"new_achievements"`
}

// AchievementWithStatus 带解锁状态的成就
type AchievementWithStatus struct {
	Achievement
	Unlocked   bool       `json:"unlocked"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
}

// UpdateGoalReq 更新每日目标请求
type UpdateGoalReq struct {
	DailyCharGoal int `json:"daily_char_goal" binding:"required,min=100,max=10000"`
}

// WritingReport 写作报告
type WritingReport struct {
	Period    string `json:"period"`     // week/month
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`

	// 本期数据
	NotesCreated int64 `json:"notes_created"`
	CharsWritten int64 `json:"chars_written"`
	ActiveDays   int   `json:"active_days"`
	GoalsMet     int   `json:"goals_met"`

	// 与上期对比
	NotesDelta      int64 `json:"notes_delta"`
	CharsDelta      int64 `json:"chars_delta"`
	ActiveDaysDelta int   `json:"active_days_delta"`

	// 本期解锁的成就
	AchievementsEarned []Achievement `json:"achievements_earned"`

	// 每日数据（用于图表）
	DailyStats []DailyStat `json:"daily_stats"`
}

// DailyStat 每日统计
type DailyStat struct {
	Date    string `json:"date"`
	Notes   int64  `json:"notes"`
	Chars   int64  `json:"chars"`
	GoalMet bool   `json:"goal_met"`
}

// ========== 预置成就数据 ==========

// DefaultAchievements 预置成就列表
var DefaultAchievements = []Achievement{
	// 笔记数量成就
	{ID: "first_note", Name: "First Steps", NameZh: "初出茅庐", Description: "Create your first note", DescriptionZh: "创建第一篇笔记", Icon: "✏️", Category: "notes", Threshold: 1, Rarity: "common"},
	{ID: "notes_10", Name: "Getting Started", NameZh: "小有成就", Description: "Create 10 notes", DescriptionZh: "创建10篇笔记", Icon: "📓", Category: "notes", Threshold: 10, Rarity: "common"},
	{ID: "notes_50", Name: "Prolific Writer", NameZh: "笔耕不辍", Description: "Create 50 notes", DescriptionZh: "创建50篇笔记", Icon: "📚", Category: "notes", Threshold: 50, Rarity: "rare"},
	{ID: "notes_100", Name: "Century Club", NameZh: "百篇达人", Description: "Create 100 notes", DescriptionZh: "创建100篇笔记", Icon: "🏆", Category: "notes", Threshold: 100, Rarity: "epic"},
	{ID: "notes_500", Name: "Master Scribe", NameZh: "著作等身", Description: "Create 500 notes", DescriptionZh: "创建500篇笔记", Icon: "👑", Category: "notes", Threshold: 500, Rarity: "legendary"},

	// 连续天数成就
	{ID: "streak_3", Name: "Warming Up", NameZh: "初露锋芒", Description: "3-day writing streak", DescriptionZh: "连续写作3天", Icon: "🔥", Category: "streak", Threshold: 3, Rarity: "common"},
	{ID: "streak_7", Name: "Week Warrior", NameZh: "周周坚持", Description: "7-day writing streak", DescriptionZh: "连续写作7天", Icon: "💪", Category: "streak", Threshold: 7, Rarity: "rare"},
	{ID: "streak_30", Name: "Monthly Master", NameZh: "月度达人", Description: "30-day writing streak", DescriptionZh: "连续写作30天", Icon: "📅", Category: "streak", Threshold: 30, Rarity: "epic"},
	{ID: "streak_100", Name: "Unstoppable", NameZh: "势不可挡", Description: "100-day writing streak", DescriptionZh: "连续写作100天", Icon: "🚀", Category: "streak", Threshold: 100, Rarity: "legendary"},

	// 字符数成就
	{ID: "chars_1k", Name: "Wordsmith", NameZh: "初级写手", Description: "Write 1,000 characters", DescriptionZh: "累计写作1000字", Icon: "🖊️", Category: "words", Threshold: 1000, Rarity: "common"},
	{ID: "chars_10k", Name: "Storyteller", NameZh: "故事大王", Description: "Write 10,000 characters", DescriptionZh: "累计写作1万字", Icon: "📖", Category: "words", Threshold: 10000, Rarity: "rare"},
	{ID: "chars_50k", Name: "Novelist", NameZh: "小说家", Description: "Write 50,000 characters", DescriptionZh: "累计写作5万字", Icon: "🪶", Category: "words", Threshold: 50000, Rarity: "epic"},
	{ID: "chars_100k", Name: "Literary Legend", NameZh: "文学巨匠", Description: "Write 100,000 characters", DescriptionZh: "累计写作10万字", Icon: "⭐", Category: "words", Threshold: 100000, Rarity: "legendary"},

	// 目标完成成就
	{ID: "goal_first", Name: "Goal Getter", NameZh: "目标达成", Description: "Complete your first daily goal", DescriptionZh: "首次完成每日目标", Icon: "🎯", Category: "goals", Threshold: 1, Rarity: "common"},
	{ID: "goal_7", Name: "Week of Goals", NameZh: "周周达标", Description: "Complete daily goal 7 times", DescriptionZh: "累计完成7次每日目标", Icon: "✅", Category: "goals", Threshold: 7, Rarity: "rare"},
	{ID: "goal_30", Name: "Goal Machine", NameZh: "目标机器", Description: "Complete daily goal 30 times", DescriptionZh: "累计完成30次每日目标", Icon: "⚡", Category: "goals", Threshold: 30, Rarity: "epic"},
}
