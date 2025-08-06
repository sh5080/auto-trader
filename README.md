# Auto Trader

해외주식 자동화 거래 시스템입니다. Go 언어로 개발되었으며, 다양한 기술적 분석 전략을 지원합니다.

## 주요 기능

- **실시간 가격 모니터링**: WebSocket을 통한 실시간 가격 데이터 수집
- **다양한 거래 전략**: 이동평균, RSI, 볼린저 밴드 등
- **리스크 관리**: 포지션 크기 제한, 일일 손실 한도, 스탑로스
- **REST API**: 전략 관리 및 모니터링을 위한 API
- **고성능**: Go의 동시성 처리로 빠른 응답 속도

## 설치 및 실행

### 1. Go 설치
```bash
# Go 1.21 이상 필요
go version
```

### 2. 의존성 설치
```bash
go mod tidy
```

### 3. 설정 파일 수정
```bash
# config/config.yaml 파일에서 API 키와 설정을 수정
cp config/config.yaml config/config.yaml.example
```

### 4. 실행
```bash
go run cmd/trader/main.go
```

## API 엔드포인트

### 상태 확인
```
GET /health
```

### 전략 관리
```
GET /strategies                    # 전략 목록 조회
POST /strategies/:id/start         # 전략 시작
POST /strategies/:id/stop          # 전략 중지
```

### 데이터 조회
```
GET /data/price/:symbol            # 현재 가격 조회
```

### 주문 관리
```
GET /orders                        # 주문 목록 조회
```

## 지원하는 전략

### 1. 이동평균 크로스오버
- 단기 이동평균이 장기 이동평균을 상향 돌파할 때 매수
- 단기 이동평균이 장기 이동평균을 하향 돌파할 때 매도

### 2. RSI 전략
- RSI가 30 이하일 때 매수 (과매도)
- RSI가 70 이상일 때 매도 (과매수)

### 3. 볼린저 밴드
- 가격이 하단 밴드에 터치할 때 매수
- 가격이 상단 밴드에 터치할 때 매도

## 리스크 관리

- **최대 포지션 크기**: 10,000 USD
- **일일 손실 한도**: 1,000 USD
- **최대 드로우다운**: 10%
- **스탑로스**: 5%

## 프로젝트 구조

```
auto-trader/
├── cmd/
│   └── trader/           # 메인 실행 파일
├── internal/
│   ├── config/          # 설정 관리
│   ├── data/            # 데이터 수집
│   ├── strategy/        # 거래 전략
│   ├── execution/       # 주문 실행
│   └── risk/            # 리스크 관리
├── pkg/
│   ├── api/             # 외부 API 클라이언트
│   └── utils/           # 유틸리티 함수
├── config/              # 설정 파일
├── logs/                # 로그 파일
├── go.mod
└── README.md
```

## 개발 환경 설정

### 1. 개발 도구 설치
```bash
# Air (핫 리로드)
go install github.com/cosmtrek/air@latest

# Delve (디버거)
go install github.com/go-delve/delve/cmd/dlv@latest
```

### 2. 개발 모드 실행
```bash
# Air를 사용한 핫 리로드
air

# 또는 직접 실행
go run cmd/trader/main.go
```

## 테스트

```bash
# 전체 테스트 실행
go test ./...

# 특정 패키지 테스트
go test ./internal/strategy

# 벤치마크 테스트
go test -bench=. ./...
```

## 배포

### Docker 사용
```bash
# Docker 이미지 빌드
docker build -t auto-trader .

# 컨테이너 실행
docker run -p 8080:8080 auto-trader
```

### 바이너리 빌드
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o auto-trader cmd/trader/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o auto-trader cmd/trader/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o auto-trader.exe cmd/trader/main.go
```

## 로깅

로그는 JSON 형식으로 출력되며, 다음과 같은 레벨을 지원합니다:
- `INFO`: 일반 정보
- `WARN`: 경고
- `ERROR`: 오류
- `DEBUG`: 디버그 정보

## 모니터링

### 메트릭 수집
- 거래 성과
- API 응답 시간
- 오류율
- 포지션 상태

### 알림 설정
- 일일 손실 한도 도달
- 스탑로스 실행
- API 오류 발생

## 라이선스

MIT License

## 기여

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 주의사항

⚠️ **이 소프트웨어는 교육 목적으로만 제공됩니다. 실제 거래에 사용하기 전에 충분한 테스트를 진행하세요.**

- 실제 거래에는 추가적인 리스크 관리가 필요합니다
- API 키와 시크릿은 안전하게 보관하세요
- 백테스팅을 통한 전략 검증을 권장합니다 