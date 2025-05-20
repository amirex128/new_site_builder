package fileitem

import (
	"database/sql/driver"
	"errors"
)

// FileItemPermissionEnum defines permission levels for file items
type FileItemPermissionEnum string

const (
	PrivatePermission FileItemPermissionEnum = "private"
	PublicPermission  FileItemPermissionEnum = "public"
)

func (e *FileItemPermissionEnum) Scan(src interface{}) error {
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
	if !FileItemPermissionEnum(b).IsValid() {
		return errors.New("unsupported type")
	}
	*e = FileItemPermissionEnum(b)
	return nil
}

func (e FileItemPermissionEnum) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("value invalid FileItemPermissionEnum")
	}
	return string(e), nil
}

// IsValid try to validate enum value on this type
func (e FileItemPermissionEnum) IsValid() bool {
	var permissionTypes = []string{
		string(PrivatePermission),
		string(PublicPermission),
	}

	for _, permissionType := range permissionTypes {
		if permissionType == string(e) {
			return true
		}
	}
	return false
}

// OperationType defines types of file operations
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

// IsValid try to validate enum value on this type
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
