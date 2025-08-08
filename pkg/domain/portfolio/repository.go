package portfolio

import (
	"auto-trader/pkg/domain/portfolio/dto"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository 포트폴리오 리포지토리 인터페이스
type Repository interface {
	// 포트폴리오 관련
	GetPortfolio(userID string) (*Portfolio, error)
	SavePortfolio(portfolio *Portfolio) error
	DeletePortfolio(userID string) error

	// 포지션 관련
	GetPositions(userID string) ([]*Position, error)
	GetPosition(userID, symbol string) (*Position, error)
	SavePosition(position *Position) error
	DeletePosition(userID, symbol string) error

	// 주식 가격 관련
	GetStockPrice(symbol string) (*StockPrice, error)
	SaveStockPrice(price *StockPrice) error
	GetStockPrices(symbols []string) ([]*StockPrice, error)

	// 거래 내역 관련
	GetTradeHistory(userID string, q dto.GetTradeHistoryQuery) ([]*TradeHistory, error)
	SaveTradeHistory(trade *TradeHistory) error

	// 포트폴리오 요약 관련
	GetPortfolioSummary(userID string) (*PortfolioSummary, error)
	SavePortfolioSummary(summary *PortfolioSummary) error
}

// ExternalDataSource 외부 데이터 소스 인터페이스 (의존성 역전)
type ExternalDataSource interface {
	// 잔고 조회
	GetBalance(accountNo string) ([]*Position, error)

	// 현재가 조회
	GetCurrentPrice(symbol string) (*StockPrice, error)
	GetCurrentPrices(symbols []string) ([]*StockPrice, error)

	// 포트폴리오 요약 조회
	GetPortfolioSummary(userID, accountNo string) (*PortfolioSummary, error)
}

// DBRepository SQL 기반 구현체
type DBRepository struct {
	db *sqlx.DB
}

// NewDBRepository 생성자
func NewDBRepository(db *sqlx.DB) Repository {
	return &DBRepository{db: db}
}

// GetPortfolio 포트폴리오 조회
func (r *DBRepository) GetPortfolio(userID string) (*Portfolio, error) {
	var p Portfolio
	query := `
        SELECT id, user_id, total_value, total_profit, profit_rate,
               last_updated, updated_at
        FROM portfolios WHERE user_id = $1
    `
	if err := r.db.Get(&p, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("포트폴리오를 찾을 수 없습니다: %s", userID)
		}
		return nil, fmt.Errorf("포트폴리오 조회 실패: %w", err)
	}
	return &p, nil
}

// SavePortfolio 포트폴리오 저장 (UPSERT by user_id)
func (r *DBRepository) SavePortfolio(portfolio *Portfolio) error {
	if portfolio.ID == "" {
		portfolio.ID = uuid.NewString()
	}
	query := `
        INSERT INTO portfolios (
            id, user_id, total_value, total_profit, profit_rate,
            last_updated, updated_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7)
        ON CONFLICT (user_id) DO UPDATE SET
            total_value = EXCLUDED.total_value,
            total_profit = EXCLUDED.total_profit,
            profit_rate = EXCLUDED.profit_rate,
            last_updated = EXCLUDED.last_updated,
            updated_at = EXCLUDED.updated_at,
    `
	_, err := r.db.Exec(query,
		portfolio.ID,
		portfolio.UserID,
		portfolio.TotalValue,
		portfolio.TotalProfit,
		portfolio.ProfitRate,
		portfolio.LastUpdated,
		portfolio.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("포트폴리오 저장 실패: %w", err)
	}
	return nil
}

// DeletePortfolio 포트폴리오 삭제
func (r *DBRepository) DeletePortfolio(userID string) error {
	_, err := r.db.Exec(`DELETE FROM portfolios WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("포트폴리오 삭제 실패: %w", err)
	}
	return nil
}

// GetPositions 보유 주식 목록 조회
func (r *DBRepository) GetPositions(userID string) ([]*Position, error) {
	var positions []*Position
	query := `
        SELECT id, user_id, symbol, company_name, quantity, average_price, current_price,
               total_value, total_profit, profit_rate, daily_profit, daily_profit_rate,
               last_updated, updated_at
        FROM positions WHERE user_id = $1 ORDER BY updated_at DESC
    `
	if err := r.db.Select(&positions, query, userID); err != nil {
		return nil, fmt.Errorf("포지션 목록 조회 실패: %w", err)
	}
	return positions, nil
}

// GetPosition 특정 보유 주식 조회
func (r *DBRepository) GetPosition(userID, symbol string) (*Position, error) {
	var p Position
	query := `
        SELECT id, user_id, symbol, company_name, quantity, average_price, current_price,
               total_value, total_profit, profit_rate, daily_profit, daily_profit_rate,
               last_updated, updated_at
        FROM positions WHERE user_id = $1 AND symbol = $2 LIMIT 1
    `
	if err := r.db.Get(&p, query, userID, symbol); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("보유 주식을 찾을 수 없습니다: %s", symbol)
		}
		return nil, fmt.Errorf("포지션 조회 실패: %w", err)
	}
	return &p, nil
}

// SavePosition 보유 주식 저장 (UPSERT by user_id+symbol)
func (r *DBRepository) SavePosition(position *Position) error {
	if position.ID == "" {
		position.ID = uuid.NewString()
	}
	query := `
        INSERT INTO positions (
            id, user_id, symbol, company_name, quantity, average_price, current_price,
            total_value, total_profit, profit_rate, daily_profit, daily_profit_rate,
            last_updated, updated_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
        ON CONFLICT (user_id, symbol) DO UPDATE SET
            company_name = EXCLUDED.company_name,
            quantity = EXCLUDED.quantity,
            average_price = EXCLUDED.average_price,
            current_price = EXCLUDED.current_price,
            total_value = EXCLUDED.total_value,
            total_profit = EXCLUDED.total_profit,
            profit_rate = EXCLUDED.profit_rate,
            daily_profit = EXCLUDED.daily_profit,
            daily_profit_rate = EXCLUDED.daily_profit_rate,
            last_updated = EXCLUDED.last_updated,
            updated_at = EXCLUDED.updated_at,
    `
	_, err := r.db.Exec(query,
		position.ID,
		position.UserID,
		position.Symbol,
		position.CompanyName,
		position.Quantity,
		position.AveragePrice,
		position.CurrentPrice,
		position.TotalValue,
		position.TotalProfit,
		position.ProfitRate,
		position.DailyProfit,
		position.DailyProfitRate,
		position.LastUpdated,
		position.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("포지션 저장 실패: %w", err)
	}
	return nil
}

// DeletePosition 보유 주식 삭제
func (r *DBRepository) DeletePosition(userID, symbol string) error {
	_, err := r.db.Exec(`DELETE FROM positions WHERE user_id = $1 AND symbol = $2`, userID, symbol)
	if err != nil {
		return fmt.Errorf("포지션 삭제 실패: %w", err)
	}
	return nil
}

// GetStockPrice 주식 가격 조회
func (r *DBRepository) GetStockPrice(symbol string) (*StockPrice, error) {
	var sp StockPrice
	query := `
        SELECT symbol, price, change, change_rate, volume, market_cap,
               high, low, open, previous_close, timestamp
        FROM stock_prices WHERE symbol = $1
    `
	if err := r.db.Get(&sp, query, symbol); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("주식 가격을 찾을 수 없습니다: %s", symbol)
		}
		return nil, fmt.Errorf("주식 가격 조회 실패: %w", err)
	}
	return &sp, nil
}

// SaveStockPrice 주식 가격 저장 (UPSERT by symbol)
func (r *DBRepository) SaveStockPrice(price *StockPrice) error {
	query := `
        INSERT INTO stock_prices (
            symbol, price, change, change_rate, volume, market_cap,
            high, low, open, previous_close, timestamp
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (symbol) DO UPDATE SET
            price = EXCLUDED.price,
            change = EXCLUDED.change,
            change_rate = EXCLUDED.change_rate,
            volume = EXCLUDED.volume,
            market_cap = EXCLUDED.market_cap,
            high = EXCLUDED.high,
            low = EXCLUDED.low,
            open = EXCLUDED.open,
            previous_close = EXCLUDED.previous_close,
            timestamp = EXCLUDED.timestamp
    `
	_, err := r.db.Exec(query,
		price.Symbol,
		price.Price,
		price.Change,
		price.ChangeRate,
		price.Volume,
		price.MarketCap,
		price.High,
		price.Low,
		price.Open,
		price.PreviousClose,
		price.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("주식 가격 저장 실패: %w", err)
	}
	return nil
}

// GetStockPrices 여러 주식 가격 조회
func (r *DBRepository) GetStockPrices(symbols []string) ([]*StockPrice, error) {
	if len(symbols) == 0 {
		return []*StockPrice{}, nil
	}
	query, args, err := sqlx.In(`
        SELECT symbol, price, change, change_rate, volume, market_cap,
               high, low, open, previous_close, timestamp
        FROM stock_prices WHERE symbol IN (?);
    `, symbols)
	if err != nil {
		return nil, fmt.Errorf("쿼리 빌드 실패: %w", err)
	}
	query = r.db.Rebind(query)

	var prices []*StockPrice
	if err := r.db.Select(&prices, query, args...); err != nil {
		return nil, fmt.Errorf("여러 주식 가격 조회 실패: %w", err)
	}
	return prices, nil
}

// GetTradeHistory 거래 내역 조회
func (r *DBRepository) GetTradeHistory(userID string, q dto.GetTradeHistoryQuery) ([]*TradeHistory, error) {
	var trades []*TradeHistory
	base := `
        SELECT id, user_id, symbol, type, quantity, price, total, fee, timestamp
        FROM trade_histories WHERE user_id = $1
    `
	args := []interface{}{userID}
	if strings.TrimSpace(q.Symbol) != "" {
		base += " AND symbol = $2"
		args = append(args, q.Symbol)
	}
	base += " ORDER BY timestamp DESC"
	if q.Limit > 0 {
		args = append(args, q.Limit)
		base += fmt.Sprintf(" LIMIT $%d", len(args))
	}
	if q.Offset > 0 {
		args = append(args, q.Offset)
		base += fmt.Sprintf(" OFFSET $%d", len(args))
	}
	if err := r.db.Select(&trades, base, args...); err != nil {
		return nil, fmt.Errorf("거래 내역 조회 실패: %w", err)
	}
	return trades, nil
}

// SaveTradeHistory 거래 내역 저장
func (r *DBRepository) SaveTradeHistory(trade *TradeHistory) error {
	if trade.ID == "" {
		trade.ID = uuid.NewString()
	}
	query := `
        INSERT INTO trade_histories (id, user_id, symbol, type, quantity, price, total, fee, timestamp)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
    `
	_, err := r.db.Exec(query,
		trade.ID,
		trade.UserID,
		trade.Symbol,
		trade.Type,
		trade.Quantity,
		trade.Price,
		trade.Total,
		trade.Fee,
		trade.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("거래 내역 저장 실패: %w", err)
	}
	return nil
}

// GetPortfolioSummary 포트폴리오 요약 조회 (간단 집계)
func (r *DBRepository) GetPortfolioSummary(userID string) (*PortfolioSummary, error) {
	positions, err := r.GetPositions(userID)
	if err != nil {
		return nil, err
	}
	sum := &PortfolioSummary{
		Positions:     make([]Position, 0, len(positions)),
		LastUpdated:   time.Now(),
		DataFreshness: "DB",
	}
	for _, p := range positions {
		sum.Positions = append(sum.Positions, *p)
	}
	return sum, nil
}

// SavePortfolioSummary 현재 미사용
func (r *DBRepository) SavePortfolioSummary(summary *PortfolioSummary) error { return nil }
