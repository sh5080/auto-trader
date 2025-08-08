package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 애플리케이션 설정 구조체
type Config struct {
	Server           ServerConfig           `mapstructure:"server"`
	Database         DatabaseConfig         `mapstructure:"database"`
	Risk             RiskConfig             `mapstructure:"risk"`
	Trading          TradingConfig          `mapstructure:"trading"`
	ProfitManagement ProfitManagementConfig `mapstructure:"profit_management"`
	KIS              KISConfig              `mapstructure:"kis"`
	JWT              JWTConfig              `mapstructure:"jwt"`
}

// ServerConfig 서버 설정
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DatabaseConfig 데이터베이스 설정
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// RiskConfig 리스크 관리 설정
type RiskConfig struct {
	MaxPositionSize    float64 `mapstructure:"max_position_size"`
	MaxDailyLoss       float64 `mapstructure:"max_daily_loss"`
	MaxDrawdown        float64 `mapstructure:"max_drawdown"`
	StopLossPercentage float64 `mapstructure:"stop_loss_percentage"`
}

// TradingConfig 트레이딩 설정
type TradingConfig struct {
	DefaultQuantity float64       `mapstructure:"default_quantity"`
	OrderTimeout    time.Duration `mapstructure:"order_timeout"`
	RetryAttempts   int           `mapstructure:"retry_attempts"`
}

// ProfitManagementConfig 수익 관리 전략 설정
type ProfitManagementConfig struct {
	ProfitTargetPercent  float64 `mapstructure:"profit_target_percent"`
	LossThresholdPercent float64 `mapstructure:"loss_threshold_percent"`
	SellPercentage       float64 `mapstructure:"sell_percentage"`
	MaxProfitThreshold   float64 `mapstructure:"max_profit_threshold"`
	MaxLossThreshold     float64 `mapstructure:"max_loss_threshold"`
	DailyLossThreshold   float64 `mapstructure:"daily_loss_threshold"`
	DailyProfitThreshold float64 `mapstructure:"daily_profit_threshold"`
	SafeBuyAmount        float64 `mapstructure:"safe_buy_amount"`
	MinBuyAmount         float64 `mapstructure:"min_buy_amount"`
	MaxBuyAmount         float64 `mapstructure:"max_buy_amount"`
}

// KISConfig 한국투자증권 API 설정
type KISConfig struct {
	AppKey      string `mapstructure:"app_key"`
	AppSecret   string `mapstructure:"app_secret"`
	BaseURL     string `mapstructure:"base_url"`
	AccessToken string `mapstructure:"access_token"`
	IsDemo      bool   `mapstructure:"is_demo"`
}

// JWTConfig JWT 설정
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

// Load 설정 로드
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 환경변수 설정
	viper.AutomaticEnv()

	// 환경변수 키 매핑 설정
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 기본값 설정
	setDefaults()

	// 환경변수에서 포트 우선 읽기
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", "8087")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "auto_trader")
	viper.SetDefault("kis.app_key", "")
	viper.SetDefault("kis.app_secret", "")
	viper.SetDefault("kis.base_url", "https://openapi.koreainvestment.com:9443")
	viper.SetDefault("kis.access_token", "")
	viper.SetDefault("kis.is_demo", true)
	viper.SetDefault("risk.max_position_size", 10000.0)
	viper.SetDefault("risk.max_daily_loss", 1000.0)
	viper.SetDefault("risk.max_drawdown", 0.1)
	viper.SetDefault("risk.stop_loss_percentage", 0.05)
	viper.SetDefault("trading.default_quantity", 100.0)
	viper.SetDefault("trading.order_timeout", "30s")
	viper.SetDefault("trading.retry_attempts", 3)

	// 수익 관리 전략 기본값
	viper.SetDefault("profit_management.profit_target_percent", 3.0)
	viper.SetDefault("profit_management.loss_threshold_percent", -3.0)
	viper.SetDefault("profit_management.sell_percentage", 90.0)
	viper.SetDefault("profit_management.max_profit_threshold", 10.0)
	viper.SetDefault("profit_management.max_loss_threshold", -10.0)
	viper.SetDefault("profit_management.daily_loss_threshold", -1.0)
	viper.SetDefault("profit_management.daily_profit_threshold", 1.0)
	viper.SetDefault("profit_management.safe_buy_amount", 1000.0)
	viper.SetDefault("profit_management.min_buy_amount", 1000.0)
	viper.SetDefault("profit_management.max_buy_amount", 10000.0)

	// JWT 기본값
	viper.SetDefault("jwt.secret", "dev-secret-change-me")
	viper.SetDefault("jwt.access_ttl", "15m")
	viper.SetDefault("jwt.refresh_ttl", "168h") // 7d
}
