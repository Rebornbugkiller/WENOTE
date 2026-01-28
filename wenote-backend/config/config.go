package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	AI        AIConfig        `mapstructure:"ai"`
	Worker    WorkerConfig    `mapstructure:"worker"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Cleanup   CleanupConfig   `mapstructure:"cleanup"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type AIConfig struct {
	Zhipu ZhipuConfig `mapstructure:"zhipu"`
}

type ZhipuConfig struct {
	APIKey     string `mapstructure:"api_key"`
	Model      string `mapstructure:"model"`
	BaseURL    string `mapstructure:"base_url"`
	DailyQuota int    `mapstructure:"daily_quota"`
	Timeout    int    `mapstructure:"timeout"`
	MaxRetries int    `mapstructure:"max_retries"`
	RetryDelay int    `mapstructure:"retry_delay"`
}

type WorkerConfig struct {
	MaxWorkers  int `mapstructure:"max_workers"`
	QueueSize   int `mapstructure:"queue_size"`
	TaskTimeout int `mapstructure:"task_timeout"`
}

type RateLimitConfig struct {
	GlobalRate  float64 `mapstructure:"global_rate"`
	GlobalBurst int     `mapstructure:"global_burst"`
	UserRate    float64 `mapstructure:"user_rate"`
	UserBurst   int     `mapstructure:"user_burst"`
	IPRate      float64 `mapstructure:"ip_rate"`
	IPBurst     int     `mapstructure:"ip_burst"`
}

type CleanupConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Days    int  `mapstructure:"days"`
}

var GlobalConfig *Config

func InitConfig() error {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile != "" {
		viper.SetConfigFile("./config/" + configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("../config")
		viper.AddConfigPath("../../config")
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 环境变量覆盖配置
	if host := os.Getenv("DB_HOST"); host != "" {
		GlobalConfig.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			GlobalConfig.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		GlobalConfig.Database.Username = user
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		GlobalConfig.Database.Password = pass
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		GlobalConfig.Database.DBName = name
	}

	// JWT环境变量覆盖
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		GlobalConfig.JWT.Secret = secret
	}
	if expire := os.Getenv("JWT_EXPIRE"); expire != "" {
		if e, err := strconv.Atoi(expire); err == nil {
			GlobalConfig.JWT.Expire = e
		}
	}

	// AI环境变量覆盖
	if apiKey := os.Getenv("ZHIPU_API_KEY"); apiKey != "" {
		GlobalConfig.AI.Zhipu.APIKey = apiKey
	}
	if model := os.Getenv("ZHIPU_MODEL"); model != "" {
		GlobalConfig.AI.Zhipu.Model = model
	}

	return nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
	)
}
