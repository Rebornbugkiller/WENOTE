package service

import (
	"errors"
	"regexp"
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
	"wenote-backend/pkg/hash"
)

var (
	ErrEmailExists     = errors.New("邮箱已被使用")
	ErrPasswordSame    = errors.New("新密码不能与旧密码相同")
	ErrConfirmMismatch = errors.New("确认信息不匹配")
	ErrInvalidAvatar   = errors.New("无效的头像样式")
	ErrInvalidEmail    = errors.New("邮箱格式不正确")
)

// 邮箱格式验证正则
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// 预设头像列表
var validAvatarStyles = map[string]bool{
	"cat": true, "dog": true, "panda": true, "fox": true,
	"frog": true, "monkey": true, "rabbit": true, "bear": true,
	"gamepad": true, "palette": true, "books": true, "music": true,
	"star": true, "rocket": true, "gem": true, "fire": true,
}

// 预设颜色列表
var validAvatarColors = map[string]bool{
	"#fbbf24": true, "#22c55e": true, "#3b82f6": true,
	"#8b5cf6": true, "#ef4444": true, "#ec4899": true,
}

// UserService 用户服务
type UserService struct {
	userRepo         *repo.UserRepo
	noteRepo         *repo.NoteRepo
	gamificationRepo *repo.GamificationRepo
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo:         repo.NewUserRepo(),
		noteRepo:         repo.NewNoteRepo(),
		gamificationRepo: repo.NewGamificationRepo(),
	}
}

// GetProfile 获取用户信息（含统计）
func (s *UserService) GetProfile(userID uint64) (*model.UserProfileResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 获取统计数据
	totalNotes, _ := s.noteRepo.CountByUserID(userID)
	totalChars, _ := s.gamificationRepo.GetTotalChars(userID)

	var currentStreak int
	gamification, _ := s.gamificationRepo.GetOrCreateUserGamification(userID)
	if gamification != nil {
		currentStreak = gamification.CurrentStreak
	}

	return &model.UserProfileResp{
		ID:            user.ID,
		Username:      user.Username,
		Nickname:      user.Nickname,
		Email:         user.Email,
		Bio:           user.Bio,
		AvatarStyle:   user.AvatarStyle,
		AvatarColor:   user.AvatarColor,
		CreatedAt:     user.CreatedAt,
		TotalNotes:    totalNotes,
		TotalChars:    totalChars,
		CurrentStreak: currentStreak,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint64, req *model.UpdateProfileReq) (*model.UserProfileResp, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 处理邮箱更新
	if req.Email != nil {
		if *req.Email == "" {
			// 允许清空邮箱
			user.Email = ""
		} else if *req.Email != user.Email {
			// 验证邮箱格式
			if !emailRegex.MatchString(*req.Email) {
				return nil, ErrInvalidEmail
			}
			// 检查邮箱唯一性
			exists, err := s.userRepo.ExistsByEmail(*req.Email, userID)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrEmailExists
			}
			user.Email = *req.Email
		}
	}

	// 验证头像样式
	if req.AvatarStyle != nil && *req.AvatarStyle != "" {
		if !validAvatarStyles[*req.AvatarStyle] {
			return nil, ErrInvalidAvatar
		}
		user.AvatarStyle = *req.AvatarStyle
	}

	// 验证头像颜色
	if req.AvatarColor != nil && *req.AvatarColor != "" {
		if !validAvatarColors[*req.AvatarColor] {
			return nil, ErrInvalidAvatar
		}
		user.AvatarColor = *req.AvatarColor
	}

	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.GetProfile(userID)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint64, req *model.ChangePasswordReq) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 验证当前密码
	if !hash.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return ErrPasswordIncorrect
	}

	// 检查新密码是否与旧密码相同
	if hash.CheckPassword(req.NewPassword, user.PasswordHash) {
		return ErrPasswordSame
	}

	// 生成新密码哈希
	newHash, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, newHash)
}

// DeleteAccount 注销账号
func (s *UserService) DeleteAccount(userID uint64, req *model.DeleteAccountReq) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 验证密码
	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		return ErrPasswordIncorrect
	}

	// 验证确认信息
	if req.Confirm != "DELETE" {
		return ErrConfirmMismatch
	}

	return s.userRepo.Delete(userID)
}
