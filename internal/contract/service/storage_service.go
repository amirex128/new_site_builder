package service

import (
	"github.com/amirex128/new_site_builder/internal/domain/enums"
	"io"
	"time"
)

// IStorageService defines the interface for storage operations
type IStorageService interface {
	// Basic operations
	CreateFileOrDirectoryIfNotExists(serverKey, bucketName, key string, permission enums.FileItemPermissionEnum, fileStream ...io.Reader) (string, error)
	CheckFileOrDirectoryIsExists(serverKey, bucketName, key string) (bool, error)
	RemoveFileOrDirectoryIfExists(serverKey, bucketName, key string) (bool, error)
	CreateBucketIfNotExists(serverKey, bucketName string) error
	DeleteBucketIfExists(serverKey, bucketName string) error

	// File operations
	DownloadFileOrDirectory(serverKey, bucketName, key string) (io.ReadCloser, error)
	GeneratePreSignedUrl(serverKey, bucketName, key string, expiry time.Duration) (string, error)
	GenerateUrl(serverKey, bucketName, key string) string

	// Advanced operations
	RenameOrMoveFileOrDirectory(serverKey, bucketName, oldKey, newKey string, currentPolicy enums.FileItemPermissionEnum) (string, error)
	CopyFileOrDirectory(serverKey, bucketName, sourceKey, destinationDirectory string, currentPolicy enums.FileItemPermissionEnum) (bool, error)

	// Permission management
	AddOrRemoveOrChangePermission(serverKey, bucketName, key string, permission enums.FileItemPermissionEnum, justRemove bool) error
}
