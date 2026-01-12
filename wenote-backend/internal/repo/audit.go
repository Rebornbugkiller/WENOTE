package repo

import "wenote-backend/internal/model"

type AuditRepo struct{}

func NewAuditRepo() *AuditRepo {
	return &AuditRepo{}
}

// Create 创建审计日志
func (r *AuditRepo) Create(log *model.AuditLog) error {
	return DB.Create(log).Error
}
