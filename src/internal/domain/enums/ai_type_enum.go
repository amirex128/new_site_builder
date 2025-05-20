package enums

import (
	"database/sql/driver"
	"errors"
)

type AiTypeEnum string

const (
	GPT35Type  AiTypeEnum = "gpt35"
	GPT4Type   AiTypeEnum = "gpt4"
	ClaudeType AiTypeEnum = "claude"
)

func (e *AiTypeEnum) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !AiTypeEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = AiTypeEnum(b)
	return nil
}

func (e AiTypeEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid AiTypeEnum")
	}
	return string(e), nil
}

func (e AiTypeEnum) IsValid() bool {
	var types = []string{
		string(GPT35Type),
		string(GPT4Type),
		string(ClaudeType),
	}
	for _, t := range types {
		if t == string(e) {
			return true
		}
	}
	return false
}
