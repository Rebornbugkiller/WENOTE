package model

import (
	"time"
)

type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Nickname     string    `gorm:"type:varchar(100)" json:"nickname"`
	Email        string    `gorm:"type:varchar(255);index" json:"email"`
	Bio          string    `gorm:"type:text" json:"bio"`
	AvatarStyle  string    `gorm:"type:varchar(50);default:'cat'" json:"avatar_style"`
	AvatarColor  string    `gorm:"type:varchar(20);default:'#fbbf24'" json:"avatar_color"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UserResp struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=50"`
}

// UpdateProfileReq 更新资料请求
type UpdateProfileReq struct {
	Nickname    *string `json:"nickname" binding:"omitempty,max=100"`
	Email       *string `json:"email" binding:"omitempty,max=255"`
	Bio         *string `json:"bio" binding:"omitempty,max=500"`
	AvatarStyle *string `json:"avatar_style" binding:"omitempty,max=50"`
	AvatarColor *string `json:"avatar_color" binding:"omitempty,max=20"`
}

// DeleteAccountReq 注销账号请求
type DeleteAccountReq struct {
	Password string `json:"password" binding:"required"`
	Confirm  string `json:"confirm" binding:"required,eq=DELETE"`
}

// UserProfileResp 用户资料响应（含统计）
type UserProfileResp struct {
	ID            uint64    `json:"id"`
	Username      string    `json:"username"`
	Nickname      string    `json:"nickname"`
	Email         string    `json:"email"`
	Bio           string    `json:"bio"`
	AvatarStyle   string    `json:"avatar_style"`
	AvatarColor   string    `json:"avatar_color"`
	CreatedAt     time.Time `json:"created_at"`
	TotalNotes    int64     `json:"total_notes"`
	TotalChars    int64     `json:"total_chars"`
	CurrentStreak int       `json:"current_streak"`
}
