package model

import (
	"time"
)

// Tag 标签模型
// 对应数据库 tags 表
// 标签用于对笔记进行分类和标记，支持多对多关联
//
// 字段说明：
//   - ID: 标签唯一标识
//   - UserID: 所属用户 ID（每个用户有自己的标签库）
//   - Name: 标签名称
//   - CreatedAt: 创建时间
//   - NoteCount: 使用此标签的笔记数量（非数据库字段）
type Tag struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Color     string    `gorm:"type:varchar(20);default:'#6B7280'" json:"color"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 非数据库字段，用于返回统计信息
	NoteCount int64 `gorm:"-" json:"note_count,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// NoteTag 笔记-标签关联表模型
// 对应数据库 note_tags 表
// 实现笔记和标签的多对多关系
//
// 表结构：
//   - note_id: 笔记 ID（联合主键）
//   - tag_id: 标签 ID（联合主键）
//
// 一条记录表示一个笔记关联了一个标签
type NoteTag struct {
	NoteID uint64 `gorm:"primaryKey" json:"note_id"` // 笔记 ID
	TagID  uint64 `gorm:"primaryKey" json:"tag_id"`  // 标签 ID
}

// TableName 指定表名
func (NoteTag) TableName() string {
	return "note_tags"
}

// ========== 请求/响应 DTO ==========

// TagCreateReq 创建标签请求
// 用于 POST /api/v1/tags
type TagCreateReq struct {
	Name  string `json:"name" binding:"required,max=100"` // 标签名称，必填，最大 100 字符
	Color string `json:"color" binding:"max=20"`          // 标签颜色，可选
}

// TagUpdateReq 更新标签请求
// 用于 PATCH /api/v1/tags/:id
type TagUpdateReq struct {
	Name  *string `json:"name" binding:"omitempty,max=100"` // 标签名称
	Color *string `json:"color" binding:"omitempty,max=20"` // 标签颜色
}

// TagListResp 标签列表响应
// 用于 GET /api/v1/tags
type TagListResp struct {
	List []*Tag `json:"list"` // 标签列表
}
