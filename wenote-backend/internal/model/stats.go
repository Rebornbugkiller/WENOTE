package model

// StatsOverview 统计概览
type StatsOverview struct {
	TotalNotes      int64 `json:"total_notes"`       // 总笔记数
	TotalNotebooks  int64 `json:"total_notebooks"`   // 总笔记本数
	TotalTags       int64 `json:"total_tags"`        // 总标签数
	ThisWeekNotes   int64 `json:"this_week_notes"`   // 本周新增笔记数
	TotalWords      int64 `json:"total_words"`       // 总字数
	ThisWeekWords   int64 `json:"this_week_words"`   // 本周新增字数
}

// TrendData 趋势数据点
type TrendData struct {
	Date  string `json:"date"`  // 日期 YYYY-MM-DD
	Count int64  `json:"count"` // 数量
}

// TagStat 标签统计
type TagStat struct {
	TagName string `json:"tag_name"` // 标签名称
	Count   int64  `json:"count"`    // 使用次数
	Color   string `json:"color"`    // 标签颜色
}

// NotebookStat 笔记本统计
type NotebookStat struct {
	NotebookName string `json:"notebook_name"` // 笔记本名称
	Count        int64  `json:"count"`         // 笔记数量
}



