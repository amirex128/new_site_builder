package fileitem

// FileItemPermissionEnum defines permission levels for file items
type FileItemPermissionEnum int

const (
	Private FileItemPermissionEnum = iota
	Public
)

// OperationType defines types of file operations
type OperationType int

const (
	Copy OperationType = iota
	Move
	Rename
)
