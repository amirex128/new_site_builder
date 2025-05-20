package enums

import (
	"database/sql/driver"
	"errors"
)

type OperationType string

const (
	CopyOperation   OperationType = "copy"
	MoveOperation   OperationType = "move"
	RenameOperation OperationType = "rename"
)

func (e *OperationType) Scan(src interface{}) error {
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
	if !OperationType(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = OperationType(b)
	return nil
}

func (e OperationType) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid OperationType")
	}
	return string(e), nil
}

func (e OperationType) IsValid() bool {
	var operationTypes = []string{
		string(CopyOperation),
		string(MoveOperation),
		string(RenameOperation),
	}
	for _, operationType := range operationTypes {
		if operationType == string(e) {
			return true
		}
	}
	return false
}
