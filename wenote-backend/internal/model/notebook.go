package model

import (
	"time"
)

// Notebook 笔记本模型
// 对应数据库 notebooks 表
// 笔记本是笔记的容器，用于组织和分类笔记
//
// 字段说明：
//   - ID: 笔记本唯一标识
//   - UserID: 所属用户 ID，建立索引加速查询
//   - Name: 笔记本名称
//   - IsDefault: 是否为默认笔记本（不可删除）
//   - CreatedAt: 创建时间
//   - UpdatedAt: 更新时间
//   - NoteCount: 笔记数量（非数据库字段，通过查询计算）
type Notebook struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	IsDefault bool      `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联字段（非数据库字段）
	// gorm:"-" 表示 GORM 不会将此字段映射到数据库
	NoteCount int64 `gorm:"-" json:"note_count"`
}

// TableName 指定表名
func (Notebook) TableName() string {
	return "notebooks"
}

// ========== 请求/响应 DTO ==========

// NotebookCreateReq 创建笔记本请求
// 用于 POST /api/v1/notebooks
type NotebookCreateReq struct {
	Name string `json:"name" binding:"required,max=255"` // 笔记本名称，必填，最大 255 字符
}

// NotebookUpdateReq 更新笔记本请求
// 用于 PATCH /api/v1/notebooks/:id
type NotebookUpdateReq struct {
	Name string `json:"name" binding:"required,max=255"` // 新的笔记本名称
}

// NotebookListResp 笔记本列表响应
// 用于 GET /api/v1/notebooks
type NotebookListResp struct {
	List []*Notebook `json:"list"` // 笔记本列表
}
