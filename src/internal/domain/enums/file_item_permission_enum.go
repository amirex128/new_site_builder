package enums

import (
	"database/sql/driver"
	"errors"
)

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
