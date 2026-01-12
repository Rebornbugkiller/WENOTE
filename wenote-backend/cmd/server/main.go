package main

import (
	"wenote-backend/config"
	"wenote-backend/internal/repo"
	"wenote-backend/internal/router"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/ai"
	"wenote-backend/pkg/alert"
	"wenote-backend/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		fmt.Println("初始化配置失败:", err)
		os.Exit(1)
	}

	logger.Init(config.GlobalConfig.Server.Mode)
	logger.Info("配置加载成功")

	alert.NewFeishuClient(alert.FeishuConfig{
		Enabled:    config.GlobalConfig.Alert.Feishu.Enabled,
		WebhookURL: config.GlobalConfig.Alert.Feishu.WebhookURL,
	})
	logger.Info("飞书告警客户端初始化成功")

	if err := repo.InitDB(); err != nil {
		logger.Error("初始化数据库失败", "error", err)
		os.Exit(1)
	}

	aiClient := ai.NewZhipuClient(ai.ZhipuConfig{
		APIKey:     config.GlobalConfig.AI.Zhipu.APIKey,
		Model:      config.GlobalConfig.AI.Zhipu.Model,
		BaseURL:    config.GlobalConfig.AI.Zhipu.BaseURL,
		Timeout:    config.GlobalConfig.AI.Zhipu.Timeout,
		MaxRetries: config.GlobalConfig.AI.Zhipu.MaxRetries,
		RetryDelay: config.GlobalConfig.AI.Zhipu.RetryDelay,
	})
	logger.Info("AI 客户端初始化成功")

	service.InitGlobalDeps(aiClient)

	noteService := service.NewNoteService()
	stopCleanup := startCleanupScheduler(noteService)

	r := router.SetupRouter()

	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	logger.Info("服务器启动中", "addr", addr)

	go func() {
		if err := r.Run(addr); err != nil {
			logger.Error("服务器启动失败", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	close(stopCleanup)

	if err := repo.CloseDB(); err != nil {
		logger.Error("关闭数据库连接失败，服务未正常关闭", "error", err)
		os.Exit(1)
	} else {
		logger.Info("服务器已关闭")
	}
}

func startCleanupScheduler(noteService *service.NoteService) chan struct{} {
	stop := make(chan struct{})
	cfg := config.GlobalConfig.Cleanup
	if !cfg.Enabled {
		logger.Info("清理任务已禁用")
		return stop
	}

	go func() {
		now := time.Now()
		todayAt2am := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
		var firstRunTime time.Time
		if now.Before(todayAt2am) {
			firstRunTime = todayAt2am
		} else {
			firstRunTime = todayAt2am.Add(24 * time.Hour)
		}

		logger.Info("清理任务已启动", "next_run", firstRunTime, "days", cfg.Days)

		waitDuration := firstRunTime.Sub(now)
		firstTimer := time.NewTimer(waitDuration)
		defer firstTimer.Stop()

		select {
		case <-stop:
			logger.Info("清理任务已停止")
			return
		case <-firstTimer.C:
			count, err := noteService.CleanupDeletedNotes(cfg.Days)
			if err != nil {
				logger.Error("清理任务失败", "error", err)
			} else {
				logger.Info("清理任务完成", "deleted_count", count)
			}
		}

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				logger.Info("清理任务已停止")
				return
			case <-ticker.C:
				count, err := noteService.CleanupDeletedNotes(cfg.Days)
				if err != nil {
					logger.Error("清理任务失败", "error", err)
				} else {
					logger.Info("清理任务完成", "deleted_count", count)
				}
			}
		}
	}()

	return stop
}
