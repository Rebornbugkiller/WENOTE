package model

import "time"

// AuditLog 审计日志
type AuditLog struct {
	ID           uint64                 `json:"id" gorm:"primaryKey"`
	UserID       uint64                 `json:"user_id" gorm:"not null;index:idx_user_action"`
	Action       string                 `json:"action" gorm:"size:50;not null;index:idx_user_action"`
	ResourceType string                 `json:"resource_type" gorm:"size:50;not null"`
	ResourceID   uint64                 `json:"resource_id"`
	Details      map[string]interface{} `json:"details" gorm:"serializer:json"`
	IPAddress    string                 `json:"ip_address" gorm:"size:50"`
	CreatedAt    time.Time              `json:"created_at" gorm:"index:idx_created"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
