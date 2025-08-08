package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// Validator DTO 검증 인터페이스
type Validator interface {
	Validate() error
}

// ValidateStruct 구조체의 필수 필드들을 자동으로 검증
// `validate:"required"` 태그가 있는 필드들을 검증
func ValidateStruct(v interface{}) error {
	var errors ValidationErrors

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("구조체가 아닙니다")
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// validate 태그 확인
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// required 검증
		if strings.Contains(validateTag, "required") {
			if err := validateRequiredField(field); err != nil {
				errors.AddError(fieldType.Name, err.Error())
			}
		}

		// min 검증
		if strings.Contains(validateTag, "min=") {
			minValue := extractTagValue(validateTag, "min=")
			if err := validateMinLength(field, minValue); err != nil {
				errors.AddError(fieldType.Name, err.Error())
			}
		}

		// max 검증
		if strings.Contains(validateTag, "max=") {
			maxValue := extractTagValue(validateTag, "max=")
			if err := validateMaxLength(field, maxValue); err != nil {
				errors.AddError(fieldType.Name, err.Error())
			}
		}

		// enum 검증
		if strings.Contains(validateTag, "enum=") {
			enumValues := extractTagValue(validateTag, "enum=")
			if err := validateEnum(field, fieldType.Name, enumValues); err != nil {
				errors.AddError(fieldType.Name, err.Error())
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// validateRequiredField 필수 필드 검증
func validateRequiredField(field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		if strings.TrimSpace(field.String()) == "" {
			return fmt.Errorf("필수 필드입니다")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() == 0 {
			return fmt.Errorf("0이 아닌 값이어야 합니다")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() == 0 {
			return fmt.Errorf("0이 아닌 값이어야 합니다")
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() == 0 {
			return fmt.Errorf("0이 아닌 값이어야 합니다")
		}
	case reflect.Ptr:
		if field.IsNil() {
			return fmt.Errorf("필수 필드입니다")
		}
	case reflect.Slice, reflect.Array:
		if field.Len() == 0 {
			return fmt.Errorf("빈 배열이 아닌 값이어야 합니다")
		}
	case reflect.Map:
		if field.Len() == 0 {
			return fmt.Errorf("빈 맵이 아닌 값이어야 합니다")
		}
	}

	return nil
}

// validateMinLength 최소 길이 검증
func validateMinLength(field reflect.Value, minStr string) error {
	min, err := parseInt(minStr)
	if err != nil {
		return fmt.Errorf("최소값 파싱 오류: %v", err)
	}

	switch field.Kind() {
	case reflect.String:
		if len(strings.TrimSpace(field.String())) < min {
			return fmt.Errorf("최소 %d자 이상이어야 합니다", min)
		}
	case reflect.Slice, reflect.Array:
		if field.Len() < min {
			return fmt.Errorf("최소 %d개 이상이어야 합니다", min)
		}
	}

	return nil
}

// validateMaxLength 최대 길이 검증
func validateMaxLength(field reflect.Value, maxStr string) error {
	max, err := parseInt(maxStr)
	if err != nil {
		return fmt.Errorf("최대값 파싱 오류: %v", err)
	}

	switch field.Kind() {
	case reflect.String:
		if len(strings.TrimSpace(field.String())) > max {
			return fmt.Errorf("최대 %d자까지 가능합니다", max)
		}
	case reflect.Slice, reflect.Array:
		if field.Len() > max {
			return fmt.Errorf("최대 %d개까지 가능합니다", max)
		}
	}

	return nil
}

// validateEnum 열거형 값 검증
func validateEnum(field reflect.Value, fieldName, enumStr string) error {
	allowedValues := strings.Split(enumStr, ",")

	switch field.Kind() {
	case reflect.String:
		value := field.String()
		for _, allowed := range allowedValues {
			if value == strings.TrimSpace(allowed) {
				return nil
			}
		}
		return fmt.Errorf("허용된 값: %s", strings.Join(allowedValues, ", "))
	}

	return nil
}

// extractTagValue 태그에서 값 추출
func extractTagValue(tag, prefix string) string {
	start := strings.Index(tag, prefix)
	if start == -1 {
		return ""
	}

	start += len(prefix)
	end := strings.Index(tag[start:], " ")
	if end == -1 {
		return tag[start:]
	}

	return tag[start : start+end]
}

// parseInt 문자열을 정수로 변환
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// ValidateMultiple 여러 검증 함수를 순차적으로 실행
func ValidateMultiple(validations ...func() error) error {
	var errors ValidationErrors

	for _, validation := range validations {
		if err := validation(); err != nil {
			if validationErr, ok := err.(ValidationError); ok {
				errors.AddError(validationErr.Field, validationErr.Message)
			} else {
				errors.AddError("unknown", err.Error())
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
