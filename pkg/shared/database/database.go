package database

import (
	"fmt"
	"time"

	"auto-trader/pkg/shared/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// DB 데이터베이스 연결 인터페이스
type DB interface {
	GetDB() *sqlx.DB
	Close() error
	Ping() error
}

// Database 데이터베이스 연결 구조체
type Database struct {
	db *sqlx.DB
}

// NewDatabase 새로운 데이터베이스 연결 생성
func NewDatabase(cfg *config.DatabaseConfig) (DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %w", err)
	}

	// 연결 풀 설정
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 연결 테스트
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("데이터베이스 핑 실패: %w", err)
	}

	logrus.Info("✅ 데이터베이스 연결 성공")

	return &Database{db: db}, nil
}

// GetDB 데이터베이스 인스턴스 반환
func (d *Database) GetDB() *sqlx.DB {
	return d.db
}

// Close 데이터베이스 연결 종료
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Ping 데이터베이스 연결 상태 확인
func (d *Database) Ping() error {
	return d.db.Ping()
}
