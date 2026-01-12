package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ========== 自定义类型 ==========

// AIStatus AI 任务状态枚举
// 用于跟踪笔记的 AI 处理进度
type AIStatus string

const (
	AIStatusPending AIStatus = "pending" // 等待处理
	AIStatusRunning AIStatus = "running" // 正在处理
	AIStatusDone    AIStatus = "done"    // 处理完成
	AIStatusFailed  AIStatus = "failed"  // 处理失败
)

// StringSlice 自定义字符串切片类型
// 用于在数据库中存储 JSON 数组（如 AI 建议的标签列表）
//
// 实现了 sql.Scanner 和 driver.Valuer 接口，
// 使 GORM 能够自动将 []string 与数据库 JSON 字段互相转换
//
// 使用示例：
//
//	type Note struct {
//	    SuggestedTags StringSlice `gorm:"type:json"`
//	}
//	note.SuggestedTags = []string{"Go", "后端", "API"}
type StringSlice []string

// Scan 实现 sql.Scanner 接口
// 从数据库读取 JSON 数据并解析为 []string
//
// 参数 value: 数据库返回的原始值（通常是 []byte）
// 返回值: 解析错误（如果有）
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// Value 实现 driver.Valuer 接口
// 将 []string 序列化为 JSON 存入数据库
//
// 返回值: JSON 字节数组和可能的错误
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// ========== 笔记模型 ==========

// Note 笔记模型
// 对应数据库 notes 表，是系统的核心数据模型
//
// 字段分组：
//  1. 基础字段：ID、UserID、NotebookID、Title、Content
//  2. AI 相关：Summary、SummaryLen、SuggestedTags、AIStatus、AIError
//  3. 状态字段：IsPinned、IsStarred
//  4. 时间字段：DeletedAt（软删除）、CreatedAt、UpdatedAt
//  5. 关联字段：Tags（多对多关联）
//
// 特性：
//   - 支持软删除（DeletedAt 不为空表示已删除）
//   - 支持 AI 自动生成摘要和标签建议
//   - 支持置顶和星标功能
//   - 支持全文搜索（MySQL ngram 索引）
type Note struct {
	// ===== 基础字段 =====
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64 `gorm:"index;not null" json:"user_id"`       // 所属用户
	NotebookID uint64 `gorm:"index;not null" json:"notebook_id"`   // 所属笔记本
	Title      string `gorm:"type:varchar(255)" json:"title"`      // 笔记标题
	Content    string `gorm:"type:longtext" json:"content"`        // 笔记内容（支持大文本）

	// ===== AI 相关字段 =====
	Summary       string      `gorm:"type:text" json:"summary,omitempty"`           // AI 生成的摘要
	SummaryLen    int         `gorm:"default:200" json:"summary_len"`               // 期望的摘要长度
	SuggestedTags StringSlice `gorm:"type:json" json:"suggested_tags,omitempty"`    // AI 建议的标签
	AIStatus      AIStatus    `gorm:"type:enum('pending','running','done','failed');default:'pending'" json:"ai_status"` // AI 处理状态
	AIError       string      `gorm:"type:text" json:"ai_error,omitempty"`          // AI 处理错误信息

	// ===== 状态字段 =====
	IsPinned  bool `gorm:"default:false" json:"is_pinned"`  // 是否置顶
	IsStarred bool `gorm:"default:false" json:"is_starred"` // 是否星标

	// ===== 软删除 =====
	// DeletedAt 不为 nil 表示笔记已被删除（在回收站中）
	// 软删除的笔记可以恢复，超过保留期限后会被永久删除
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// ===== 关联 =====
	// 多对多关联，通过 note_tags 中间表实现
	Tags []Tag `gorm:"many2many:note_tags" json:"tags,omitempty"`
}

// TableName 指定表名
func (Note) TableName() string {
	return "notes"
}

// ========== 请求/响应 DTO ==========

// NoteCreateReq 创建笔记请求
// 用于 POST /api/v1/notes
type NoteCreateReq struct {
	NotebookID uint64   `json:"notebook_id" binding:"required"` // 所属笔记本 ID，必填
	Title      string   `json:"title" binding:"max=255"`        // 标题，最大 255 字符
	Content    string   `json:"content"`                        // 内容
	SummaryLen int      `json:"summary_len"`                    // 摘要长度，可选，默认 200
	TagIDs     []uint64 `json:"tag_ids"`                        // 关联的标签 ID 列表
}

// NoteUpdateReq 更新笔记请求
// 用于 PATCH /api/v1/notes/:id
// 所有字段都是可选的，使用指针类型区分"未传"和"传空值"
//
// 示例：
//   - {"title": "新标题"} - 只更新标题
//   - {"is_pinned": true} - 只更新置顶状态
//   - {"content": ""} - 清空内容
type NoteUpdateReq struct {
	Title      *string  `json:"title"`       // 新标题
	Content    *string  `json:"content"`     // 新内容
	NotebookID *uint64  `json:"notebook_id"` // 移动到新笔记本
	SummaryLen *int     `json:"summary_len"` // 新的摘要长度
	IsPinned   *bool    `json:"is_pinned"`   // 置顶状态
	IsStarred  *bool    `json:"is_starred"`  // 星标状态
	TagIDs     []uint64 `json:"tag_ids"`     // 新的标签列表（会替换原有标签）
}

// NoteListReq 笔记列表查询请求
// 用于 GET /api/v1/notes
// 支持多种筛选条件和分页
type NoteListReq struct {
	NotebookID *uint64 `form:"notebook_id"` // 按笔记本筛选
	TagID      *uint64 `form:"tag_id"`      // 按标签筛选
	IsStarred  *bool   `form:"is_starred"`  // 只看星标
	IsPinned   *bool   `form:"is_pinned"`   // 只看置顶
	Keyword    string  `form:"keyword"`     // 关键词搜索（全文搜索）
	Page       int     `form:"page,default=1"`      // 页码，默认 1
	PageSize   int     `form:"page_size,default=20"` // 每页数量，默认 20
}

// NoteListResp 笔记列表响应
type NoteListResp struct {
	Total int64   `json:"total"` // 总数
	List  []*Note `json:"list"`  // 笔记列表
	Page  int     `json:"page"`  // 当前页码
	Size  int     `json:"size"`  // 每页数量
}

// NoteTagsReq 修改笔记标签请求
// 用于 PUT /api/v1/notes/:id/tags
type NoteTagsReq struct {
	TagIDs []uint64 `json:"tag_ids" binding:"required"` // 新的标签 ID 列表
}

// ========== 批量操作请求 ==========

// BatchDeleteReq 批量删除请求
// 用于 POST /api/v1/notes/batch/delete
type BatchDeleteReq struct {
	NoteIDs []uint64 `json:"note_ids" binding:"required,min=1,max=100"` // 要删除的笔记 ID 列表
}

// BatchRestoreReq 批量恢复请求
// 用于 POST /api/v1/notes/batch/restore
type BatchRestoreReq struct {
	NoteIDs []uint64 `json:"note_ids" binding:"required,min=1,max=100"` // 要恢复的笔记 ID 列表
}

// BatchMoveReq 批量移动请求
// 用于 POST /api/v1/notes/batch/move
type BatchMoveReq struct {
	NoteIDs    []uint64 `json:"note_ids" binding:"required,min=1,max=100"` // 要移动的笔记 ID 列表
	NotebookID uint64   `json:"notebook_id" binding:"required"`            // 目标笔记本 ID
}
