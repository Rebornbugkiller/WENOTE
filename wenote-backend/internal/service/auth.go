package service

import (
	"wenote-backend/config"
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
	"wenote-backend/pkg/hash"
	"wenote-backend/pkg/jwt"
	"wenote-backend/pkg/logger"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUsernameExists    = errors.New("用户名已存在")
	ErrPasswordIncorrect = errors.New("密码错误")
	ErrUserCreateFailed  = errors.New("用户创建失败")
)

type LoginAttempt struct {
	Count       int
	LockedUntil time.Time
}

var loginAttempts = &sync.Map{}

type AuthService struct {
	userRepo   *repo.UserRepo
	jwtManager *jwt.JWTManager
}

func NewAuthService() *AuthService {
	cfg := config.GlobalConfig.JWT
	return &AuthService{
		userRepo:   repo.NewUserRepo(),
		jwtManager: jwt.NewJWTManager(cfg.Secret, cfg.Expire),
	}
}

func (s *AuthService) Register(req *model.RegisterReq) (*model.User, error) {
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, ErrUserCreateFailed
	}

	return user, nil
}

func (s *AuthService) Login(req *model.LoginReq) (*model.LoginResp, error) {
	attemptVal, found := loginAttempts.Load(req.Username)
	if found {
		attempt := attemptVal.(*LoginAttempt)
		if time.Now().Before(attempt.LockedUntil) {
			remainingSeconds := int(time.Until(attempt.LockedUntil).Seconds())
			return nil, fmt.Errorf("账号已锁定，请 %d 秒后重试", remainingSeconds)
		}
	}

	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		s.recordLoginFailure(req.Username)
		return nil, ErrUserNotFound
	}

	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		s.recordLoginFailure(req.Username)
		return nil, ErrPasswordIncorrect
	}

	loginAttempts.Delete(req.Username)

	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &model.LoginResp{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) recordLoginFailure(username string) {
	now := time.Now()

	val, _ := loginAttempts.LoadOrStore(username, &LoginAttempt{})
	att := val.(*LoginAttempt)

	att.Count++

	if att.Count >= 5 {
		att.LockedUntil = now.Add(15 * time.Minute)
		logger.Warn("账号已锁定", "username", username, "locked_until", att.LockedUntil)
	}
}

func (s *AuthService) RefreshToken(userID uint64, username string) (string, error) {
	return s.jwtManager.GenerateToken(userID, username)
}
