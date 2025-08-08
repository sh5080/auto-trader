package strategy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Repository 전략 리포지토리 인터페이스
type Repository interface {
	GetAll() ([]*StrategyDetails, error)
	GetByID(id string) (*StrategyDetails, error)
	GetByUserID(userID string) ([]*StrategyDetails, error)
	Create(strategy *StrategyDetails) error
	Update(strategy *StrategyDetails) error
	Delete(id string) error
	GetConfig(id string) (*StrategyConfig, error)
	SaveConfig(config *StrategyConfig) error
	GetStatus(id string) (*StrategyStatus, error)
	SaveStatus(status *StrategyStatus) error
	GetPerformance(id string) (*StrategyPerformance, error)
	SavePerformance(performance *StrategyPerformance) error
}

type DBRepository struct {
	db *sqlx.DB
}

func NewDBRepository(db *sqlx.DB) Repository {
	return &DBRepository{
		db: db,
	}
}

// UserStrategy 데이터베이스 모델
type UserStrategy struct {
	ID          string          `db:"id"`
	UserID      string          `db:"user_id"`
	StrategyID  string          `db:"strategy_id"`
	TemplateID  sql.NullString  `db:"template_id"`
	Name        string          `db:"name"`
	Description sql.NullString  `db:"description"`
	Symbol      string          `db:"symbol"`
	UserInputs  json.RawMessage `db:"user_inputs"`
	Settings    json.RawMessage `db:"settings"`
	Active      bool            `db:"active"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

// GetAll 모든 전략 조회
func (r *DBRepository) GetAll() ([]*StrategyDetails, error) {
	query := `
		SELECT 
			us.id, us.user_id, us.strategy_id, us.template_id, us.name, 
			us.description, us.symbol, us.user_inputs, us.settings, 
			us.active, us.created_at, us.updated_at,
			COALESCE(ss.status, 'inactive') as status
		FROM user_strategies us
		LEFT JOIN strategy_status ss ON us.id = ss.strategy_id
		ORDER BY us.created_at DESC
	`

	var userStrategies []struct {
		UserStrategy
		Status string `db:"status"`
	}

	err := r.db.Select(&userStrategies, query)
	if err != nil {
		return nil, fmt.Errorf("전략 목록 조회 실패: %w", err)
	}

	var strategies []*StrategyDetails
	for _, us := range userStrategies {
		strategy, err := r.convertToStrategyDetails(&us.UserStrategy, us.Status)
		if err != nil {
			logrus.Errorf("전략 변환 실패 (ID: %s): %v", us.ID, err)
			continue
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// GetByID ID로 전략 조회
func (r *DBRepository) GetByID(id string) (*StrategyDetails, error) {
	query := `
		SELECT 
			us.id, us.user_id, us.strategy_id, us.template_id, us.name,
			us.description, us.symbol, us.user_inputs, us.settings,
			us.active, us.created_at, us.updated_at,
			COALESCE(ss.status, 'inactive') as status
		FROM user_strategies us
		LEFT JOIN strategy_status ss ON us.id = ss.strategy_id
		WHERE us.id = $1
	`

	var result struct {
		UserStrategy
		Status string `db:"status"`
	}

	err := r.db.Get(&result, query, id)
	if err != nil {
		return nil, fmt.Errorf("전략 조회 실패: %w", err)
	}

	return r.convertToStrategyDetails(&result.UserStrategy, result.Status)
}

// GetByUserID 사용자별 전략 조회
func (r *DBRepository) GetByUserID(userID string) ([]*StrategyDetails, error) {
	query := `
		SELECT 
			us.id, us.user_id, us.strategy_id, us.template_id, us.name,
			us.description, us.symbol, us.user_inputs, us.settings,
			us.active, us.created_at, us.updated_at,
			COALESCE(ss.status, 'inactive') as status
		FROM user_strategies us
		LEFT JOIN strategy_status ss ON us.id = ss.strategy_id
		WHERE us.user_id = $1
		ORDER BY us.created_at DESC
	`

	var userStrategies []struct {
		UserStrategy
		Status string `db:"status"`
	}

	err := r.db.Select(&userStrategies, query, userID)
	if err != nil {
		return nil, fmt.Errorf("사용자 전략 조회 실패: %w", err)
	}

	var strategies []*StrategyDetails
	for _, us := range userStrategies {
		strategy, err := r.convertToStrategyDetails(&us.UserStrategy, us.Status)
		if err != nil {
			logrus.Errorf("전략 변환 실패 (ID: %s): %v", us.ID, err)
			continue
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}

// Create 새로운 전략 생성
func (r *DBRepository) Create(strategy *StrategyDetails) error {
	// UUID 생성
	if strategy.ID == "" {
		strategy.ID = uuid.New().String()
	}

	// JSON 직렬화
	userInputsJSON, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("user_inputs 직렬화 실패: %w", err)
	}

	settingsJSON, err := json.Marshal(map[string]interface{}{
		"quantity":      100,
		"max_positions": 1,
	})
	if err != nil {
		return fmt.Errorf("settings 직렬화 실패: %w", err)
	}

	query := `
		INSERT INTO user_strategies 
		(id, user_id, strategy_id, name, description, symbol, user_inputs, settings, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(query,
		strategy.ID,
		"system", // TODO: 실제 사용자 ID
		strategy.ID,
		strategy.Name,
		strategy.Description,
		strategy.Symbols[0], // 첫 번째 심볼 사용
		userInputsJSON,
		settingsJSON,
		strategy.Active,
	)

	if err != nil {
		return fmt.Errorf("전략 생성 실패: %w", err)
	}

	// 상태 초기화
	return r.SaveStatus(&StrategyStatus{
		ID:             strategy.ID,
		Status:         "inactive",
		LastExecution:  time.Time{},
		ExecutionCount: 0,
		Uptime:         0,
	})
}

// Update 전략 업데이트
func (r *DBRepository) Update(strategy *StrategyDetails) error {
	query := `
		UPDATE user_strategies 
		SET name = $2, description = $3, active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query,
		strategy.ID,
		strategy.Name,
		strategy.Description,
		strategy.Active,
	)

	if err != nil {
		return fmt.Errorf("전략 업데이트 실패: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("업데이트 결과 확인 실패: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("전략을 찾을 수 없습니다: %s", strategy.ID)
	}

	return nil
}

// Delete 전략 삭제
func (r *DBRepository) Delete(id string) error {
	// 트랜잭션 시작
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 실패: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 연관 데이터는 CASCADE로 자동 삭제됨
	query := `DELETE FROM user_strategies WHERE id = $1`

	result, err := tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("전략 삭제 실패: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("삭제 결과 확인 실패: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("전략을 찾을 수 없습니다: %s", id)
	}

	return tx.Commit()
}

// GetConfig 전략 설정 조회
func (r *DBRepository) GetConfig(id string) (*StrategyConfig, error) {
	// 현재는 설정이 user_strategies 테이블에 포함됨
	// TODO: 별도 설정 테이블 구현 시 수정
	return &StrategyConfig{
		ID:         id,
		Parameters: make(map[string]interface{}),
		Enabled:    true,
		Schedule:   "",
		RiskLimits: RiskLimits{},
	}, nil
}

// SaveConfig 전략 설정 저장
func (r *DBRepository) SaveConfig(config *StrategyConfig) error {
	// TODO: 별도 설정 테이블 구현
	return nil
}

// GetStatus 전략 상태 조회
func (r *DBRepository) GetStatus(id string) (*StrategyStatus, error) {
	query := `
		SELECT strategy_id, status, last_execution, execution_count, 
		       error_message, uptime_seconds
		FROM strategy_status 
		WHERE strategy_id = $1
	`

	var status StrategyStatus
	var uptimeSeconds sql.NullInt64
	var lastExecution sql.NullTime
	var errorMessage sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&status.ID,
		&status.Status,
		&lastExecution,
		&status.ExecutionCount,
		&errorMessage,
		&uptimeSeconds,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// 기본 상태 반환
			return &StrategyStatus{
				ID:             id,
				Status:         "inactive",
				LastExecution:  time.Time{},
				ExecutionCount: 0,
				Uptime:         0,
			}, nil
		}
		return nil, fmt.Errorf("전략 상태 조회 실패: %w", err)
	}

	if lastExecution.Valid {
		status.LastExecution = lastExecution.Time
	}
	if errorMessage.Valid {
		status.ErrorMessage = errorMessage.String
	}
	if uptimeSeconds.Valid {
		status.Uptime = uptimeSeconds.Int64
	}

	return &status, nil
}

// SaveStatus 전략 상태 저장
func (r *DBRepository) SaveStatus(status *StrategyStatus) error {
	query := `
		INSERT INTO strategy_status 
		(strategy_id, status, last_execution, execution_count, error_message, uptime_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (strategy_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			last_execution = EXCLUDED.last_execution,
			execution_count = EXCLUDED.execution_count,
			error_message = EXCLUDED.error_message,
			uptime_seconds = EXCLUDED.uptime_seconds,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query,
		status.ID,
		status.Status,
		status.LastExecution,
		status.ExecutionCount,
		status.ErrorMessage,
		status.Uptime,
	)

	if err != nil {
		return fmt.Errorf("전략 상태 저장 실패: %w", err)
	}

	return nil
}

// GetPerformance 전략 성과 조회
func (r *DBRepository) GetPerformance(id string) (*StrategyPerformance, error) {
	query := `
		SELECT strategy_id, total_return, win_rate, profit_loss, trade_count,
		       last_trade_time, max_drawdown, sharpe_ratio
		FROM strategy_performance 
		WHERE strategy_id = $1
	`

	var perf StrategyPerformance
	var lastTradeTime sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&perf.StrategyID,
		&perf.TotalReturn,
		&perf.WinRate,
		&perf.ProfitLoss,
		&perf.TradeCount,
		&lastTradeTime,
		&perf.MaxDrawdown,
		&perf.SharpeRatio,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// 기본 성과 반환
			return &StrategyPerformance{
				StrategyID:    id,
				TotalReturn:   0.0,
				WinRate:       0.0,
				ProfitLoss:    0.0,
				TradeCount:    0,
				LastTradeTime: time.Time{},
				MaxDrawdown:   0.0,
				SharpeRatio:   0.0,
			}, nil
		}
		return nil, fmt.Errorf("전략 성과 조회 실패: %w", err)
	}

	if lastTradeTime.Valid {
		perf.LastTradeTime = lastTradeTime.Time
	}

	return &perf, nil
}

// SavePerformance 전략 성과 저장
func (r *DBRepository) SavePerformance(performance *StrategyPerformance) error {
	query := `
		INSERT INTO strategy_performance 
		(strategy_id, total_return, win_rate, profit_loss, trade_count, 
		 last_trade_time, max_drawdown, sharpe_ratio)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (strategy_id) 
		DO UPDATE SET 
			total_return = EXCLUDED.total_return,
			win_rate = EXCLUDED.win_rate,
			profit_loss = EXCLUDED.profit_loss,
			trade_count = EXCLUDED.trade_count,
			last_trade_time = EXCLUDED.last_trade_time,
			max_drawdown = EXCLUDED.max_drawdown,
			sharpe_ratio = EXCLUDED.sharpe_ratio,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query,
		performance.StrategyID,
		performance.TotalReturn,
		performance.WinRate,
		performance.ProfitLoss,
		performance.TradeCount,
		performance.LastTradeTime,
		performance.MaxDrawdown,
		performance.SharpeRatio,
	)

	if err != nil {
		return fmt.Errorf("전략 성과 저장 실패: %w", err)
	}

	return nil
}

// ==================== 헬퍼 메서드들 ====================

// convertToStrategyDetails UserStrategy를 StrategyDetails로 변환
func (r *DBRepository) convertToStrategyDetails(us *UserStrategy, status string) (*StrategyDetails, error) {
	return &StrategyDetails{
		ID:          us.ID,
		Name:        us.Name,
		Description: us.Description.String,
		Active:      us.Active,
		Symbols:     []string{us.Symbol}, // 단일 심볼을 배열로 변환
		CreatedAt:   us.CreatedAt,
		UpdatedAt:   us.UpdatedAt,
		Status:      status,
		Performance: nil, // 필요 시 별도 조회
	}, nil
}
