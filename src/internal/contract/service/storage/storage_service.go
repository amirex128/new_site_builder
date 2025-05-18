package storage

import (
	"io"
	"time"
)

// IStorageService defines the interface for storage operations
type IStorageService interface {
	// Basic operations
	CreateBucketIfNotExists(serverKey, bucketName string) error
	CreateFileOrDirectoryIfNotExists(serverKey, bucketName, path string, permission int, content ...io.Reader) error
	CheckFileOrDirectoryIsExists(serverKey, bucketName, path string) (bool, error)

	// File operations
	RenameOrMoveFileOrDirectory(serverKey, bucketName, sourcePath, destinationPath string, permission int) error
	CopyFileOrDirectory(serverKey, bucketName, sourcePath, destinationPath string, permission int) error
	RemoveFileOrDirectoryIfExists(serverKey, bucketName, path string) error

	// Access management
	AddOrRemoveOrChangePermission(serverKey, bucketName, path string, permission int) error

	// URL generation
	GenerateUrl(serverKey, bucketName, path string) string
	GeneratePreSignedUrl(serverKey, bucketName, path string, expiry time.Duration) (string, error)

	// Download
	DownloadFileOrDirectory(serverKey, bucketName, path string) (io.Reader, error)
}
