#!/bin/bash

# Swagger 문서 생성 스크립트

# 사용법 출력 함수
show_usage() {
    echo "📚 Swagger 문서 생성 스크립트"
    echo ""
    echo "사용법:"
    echo "  $0 [옵션]"
    echo ""
    echo "옵션:"
    echo "  -h, --help              도움말 출력"
    echo "  -e, --exclude DIRS      제외할 디렉토리 (쉼표로 구분)"
    echo "  -i, --include DIRS      포함할 디렉토리만 (쉼표로 구분)"
    echo "  -o, --output DIR        출력 디렉토리 (기본값: docs)"
    echo ""
    echo "환경 변수:"
    echo "  SWAGGER_EXCLUDE_DIRS    제외할 디렉토리 (쉼표로 구분)"
    echo "  SWAGGER_INCLUDE_DIRS    포함할 디렉토리만 (쉼표로 구분)"
    echo "  SWAGGER_OUTPUT_DIR      출력 디렉토리"
    echo ""
    echo "예시:"
    echo "  $0                                    # 모든 도메인 포함"
    echo "  $0 -e test,mock                       # test, mock 디렉토리 제외"
    echo "  $0 -i strategy,portfolio              # strategy, portfolio만 포함"
    echo "  $0 -o custom_docs                     # custom_docs 디렉토리에 출력"
    echo ""
}

# 명령행 인수 파싱
OUTPUT_DIR="docs"
EXCLUDE_DIRS=("test" "mock" "internal")  # 기본 제외 디렉토리
INCLUDE_ONLY=""  # 특정 디렉토리만 포함 (비어있으면 모든 디렉토리 포함)

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -e|--exclude)
            IFS=',' read -ra EXCLUDE_DIRS <<< "$2"
            shift 2
            ;;
        -i|--include)
            INCLUDE_ONLY="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        *)
            echo "❌ 알 수 없는 옵션: $1"
            show_usage
            exit 1
            ;;
    esac
done

# 환경 변수로 설정된 값이 있으면 덮어쓰기
if [ -n "$SWAGGER_EXCLUDE_DIRS" ]; then
    IFS=',' read -ra EXCLUDE_DIRS <<< "$SWAGGER_EXCLUDE_DIRS"
fi

if [ -n "$SWAGGER_INCLUDE_DIRS" ]; then
    INCLUDE_ONLY="$SWAGGER_INCLUDE_DIRS"
fi

if [ -n "$SWAGGER_OUTPUT_DIR" ]; then
    OUTPUT_DIR="$SWAGGER_OUTPUT_DIR"
fi

echo "📚 Swagger 문서 생성 시작..."
echo "⚙️  설정:"
echo "  - 출력 디렉토리: $OUTPUT_DIR"
echo "  - 제외 디렉토리: ${EXCLUDE_DIRS[*]}"
if [ -n "$INCLUDE_ONLY" ]; then
    echo "  - 포함 디렉토리만: $INCLUDE_ONLY"
fi
echo ""

# swag CLI 경로 설정
SWAG_PATH=""
if command -v swag &> /dev/null; then
    SWAG_PATH="swag"
elif [ -f "$(go env GOPATH)/bin/swag" ]; then
    SWAG_PATH="$(go env GOPATH)/bin/swag"
else
    echo "❌ swag CLI가 설치되지 않았습니다."
    echo "다음 명령어로 설치하세요:"
    echo "go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

echo "🔧 swag 경로: $SWAG_PATH"

# 기존 docs 디렉토리 삭제 후 새로 생성
echo "🧹 기존 Swagger 문서 삭제 중..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"
echo "📁 $OUTPUT_DIR 디렉토리 생성 완료"

# pkg/domain 하위의 모든 디렉토리 동적 스캔
echo "🔍 pkg/domain 디렉토리 스캔 중..."
DOMAIN_DIRS="cmd/trader"

# pkg/domain 디렉토리가 존재하는지 확인
if [ -d "pkg/domain" ]; then
    # pkg/domain 하위의 모든 디렉토리를 찾아서 쉼표로 구분된 문자열로 만듦
    for dir in pkg/domain/*/; do
        if [ -d "$dir" ]; then
            dir_name=$(basename "$dir")
            
            # 제외할 디렉토리인지 확인
            should_exclude=false
            for exclude_dir in "${EXCLUDE_DIRS[@]}"; do
                if [ "$dir_name" = "$exclude_dir" ]; then
                    should_exclude=true
                    break
                fi
            done
            
            # 특정 디렉토리만 포함하는 경우
            if [ -n "$INCLUDE_ONLY" ]; then
                should_include=false
                for include_dir in $INCLUDE_ONLY; do
                    if [ "$dir_name" = "$include_dir" ]; then
                        should_include=true
                        break
                    fi
                done
                if [ "$should_include" = false ]; then
                    echo "⏭️  제외됨 (INCLUDE_ONLY): $dir_name"
                    continue
                fi
            fi
            
            # 제외할 디렉토리가 아닌 경우에만 추가
            if [ "$should_exclude" = false ]; then
                echo "📂 발견된 도메인: $dir_name"
                DOMAIN_DIRS="$DOMAIN_DIRS,$dir"
            else
                echo "⏭️  제외됨 (EXCLUDE_DIRS): $dir_name"
            fi
        fi
    done
else
    echo "⚠️  pkg/domain 디렉토리가 존재하지 않습니다."
fi

echo "📋 스캔된 디렉토리: $DOMAIN_DIRS"

# Swagger 문서 생성
echo "🔄 Swagger 문서 생성 중..."
$SWAG_PATH init \
    --dir "$DOMAIN_DIRS" \
    --output "$OUTPUT_DIR" \
    --generalInfo main.go \
    --propertyStrategy snakecase \
    --parseDependency \
    --parseInternal

if [ $? -eq 0 ]; then
    echo "✅ Swagger 문서 생성 완료!"
    echo "📁 생성된 파일:"
    ls -la "$OUTPUT_DIR"/
    echo ""
    echo "🌐 Swagger UI 접속: http://localhost:8080/swagger/"
    echo ""
    echo "📊 스캔된 도메인 디렉토리:"
    echo "$DOMAIN_DIRS" | tr ',' '\n' | grep "pkg/domain" | sed 's|pkg/domain/||' | sed 's|/||' || true
    echo ""
    echo "⚙️  설정 정보:"
    echo "  - 출력 디렉토리: $OUTPUT_DIR"
    echo "  - 제외된 디렉토리: ${EXCLUDE_DIRS[*]}"
    if [ -n "$INCLUDE_ONLY" ]; then
        echo "  - 포함된 디렉토리만: $INCLUDE_ONLY"
    fi
else
    echo "❌ Swagger 문서 생성 실패!"
    exit 1
fi 