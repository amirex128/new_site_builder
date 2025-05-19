package storage

import (
	"io"
	"time"
)

// FileItemPermissionEnum represents the permission level for a file item
type FileItemPermissionEnum int

const (
	Private FileItemPermissionEnum = iota
	Public
)

// IStorageService defines the interface for storage operations
type IStorageService interface {
	// Basic operations
	CreateFileOrDirectoryIfNotExists(serverKey, bucketName, key string, permission FileItemPermissionEnum, fileStream ...io.Reader) (string, error)
	CheckFileOrDirectoryIsExists(serverKey, bucketName, key string) (bool, error)
	RemoveFileOrDirectoryIfExists(serverKey, bucketName, key string) (bool, error)
	CreateBucketIfNotExists(serverKey, bucketName string) error
	DeleteBucketIfExists(serverKey, bucketName string) error

	// File operations
	DownloadFileOrDirectory(serverKey, bucketName, key string) (io.ReadCloser, error)
	GeneratePreSignedUrl(serverKey, bucketName, key string, expiry time.Duration) (string, error)
	GenerateUrl(serverKey, bucketName, key string) string

	// Advanced operations
	RenameOrMoveFileOrDirectory(serverKey, bucketName, oldKey, newKey string, currentPolicy FileItemPermissionEnum) (string, error)
	CopyFileOrDirectory(serverKey, bucketName, sourceKey, destinationDirectory string, currentPolicy FileItemPermissionEnum) (bool, error)

	// Permission management
	AddOrRemoveOrChangePermission(serverKey, bucketName, key string, permission FileItemPermissionEnum, justRemove bool) error
}
