package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	contractStorage "github.com/amirex128/new_site_builder/src/internal/contract/service/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorageService implements the IStorageService interface using Minio
type MinioStorageService struct {
	clients map[string]*minio.Client
	config  MinioConfig
}

// MinioConfig holds configuration for Minio storage service
type MinioConfig struct {
	Endpoints map[string]string // Map of server keys to endpoints
	AccessKey string
	SecretKey string
	UseSSL    bool
}

// NewMinioStorageService creates a new MinioStorageService instance
func NewMinioStorageService(config MinioConfig) (*MinioStorageService, error) {
	clients := make(map[string]*minio.Client)

	for serverKey, endpoint := range config.Endpoints {
		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
			Secure: config.UseSSL,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize minio client for %s: %v", serverKey, err)
		}
		clients[serverKey] = client
	}

	return &MinioStorageService{
		clients: clients,
		config:  config,
	}, nil
}

// CreateBucketIfNotExists creates a bucket if it doesn't already exist
func (s *MinioStorageService) CreateBucketIfNotExists(serverKey, bucketName string) error {
	client, exists := s.clients[serverKey]
	if !exists {
		return fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking if bucket exists: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("error creating bucket: %v", err)
		}
	}

	return nil
}

// CreateFileOrDirectoryIfNotExists creates a file or directory if it doesn't exist
func (s *MinioStorageService) CreateFileOrDirectoryIfNotExists(serverKey, bucketName, path string, permission int, content ...io.Reader) error {
	client, exists := s.clients[serverKey]
	if !exists {
		return fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()

	// Create bucket if it doesn't exist
	if err := s.CreateBucketIfNotExists(serverKey, bucketName); err != nil {
		return err
	}

	// Check if object exists
	_, err := client.StatObject(ctx, bucketName, path, minio.StatObjectOptions{})
	if err == nil {
		// Object already exists
		return nil
	}

	// For directories (in S3, directories are represented by objects with trailing slash)
	if path[len(path)-1] == '/' {
		// Create empty object with trailing slash to represent directory
		_, err = client.PutObject(ctx, bucketName, path, bytes.NewReader([]byte{}), 0, minio.PutObjectOptions{
			ContentType: "application/directory",
		})
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	} else if len(content) > 0 {
		// For files with content
		contentType := "application/octet-stream"
		if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg" {
			contentType = "image/jpeg"
		} else if filepath.Ext(path) == ".png" {
			contentType = "image/png"
		} else if filepath.Ext(path) == ".pdf" {
			contentType = "application/pdf"
		}

		_, err = client.PutObject(ctx, bucketName, path, content[0], -1, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
	} else {
		// For empty files
		_, err = client.PutObject(ctx, bucketName, path, bytes.NewReader([]byte{}), 0, minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("error creating empty file: %v", err)
		}
	}

	// Set permissions
	if permission == 1 { // Public
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": ["s3:GetObject"],
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Resource": ["arn:aws:s3:::%s/%s"],
					"Sid": ""
				}
			]
		}`, bucketName, path)

		err = client.SetObjectLockConfig(ctx, bucketName, &minio.ObjectLockConfig{
			ObjectLockEnabled: "Enabled",
		})
		if err != nil {
			return fmt.Errorf("error setting object lock config: %v", err)
		}
	}

	return nil
}

// CheckFileOrDirectoryIsExists checks if a file or directory exists
func (s *MinioStorageService) CheckFileOrDirectoryIsExists(serverKey, bucketName, path string) (bool, error) {
	client, exists := s.clients[serverKey]
	if !exists {
		return false, fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()
	_, err := client.StatObject(ctx, bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// RenameOrMoveFileOrDirectory renames or moves a file or directory
func (s *MinioStorageService) RenameOrMoveFileOrDirectory(serverKey, bucketName, sourcePath, destinationPath string, permission int) error {
	client, exists := s.clients[serverKey]
	if !exists {
		return fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()

	// In S3/Minio, moving an object means copying and then deleting the original
	src := minio.CopySource{
		Bucket: bucketName,
		Object: sourcePath,
	}

	_, err := client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: bucketName,
		Object: destinationPath,
	}, src)
	if err != nil {
		return fmt.Errorf("error copying object: %v", err)
	}

	// Delete the original
	err = client.RemoveObject(ctx, bucketName, sourcePath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("error removing original object: %v", err)
	}

	// Set permissions
	if permission == 1 { // Public
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": ["s3:GetObject"],
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Resource": ["arn:aws:s3:::%s/%s"],
					"Sid": ""
				}
			]
		}`, bucketName, destinationPath)

		err = client.SetObjectLockConfig(ctx, bucketName, &minio.ObjectLockConfig{
			ObjectLockEnabled: "Enabled",
		})
		if err != nil {
			return fmt.Errorf("error setting object lock config: %v", err)
		}
	}

	return nil
}

// CopyFileOrDirectory copies a file or directory
func (s *MinioStorageService) CopyFileOrDirectory(serverKey, bucketName, sourcePath, destinationPath string, permission int) error {
	client, exists := s.clients[serverKey]
	if !exists {
		return fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()

	src := minio.CopySource{
		Bucket: bucketName,
		Object: sourcePath,
	}

	_, err := client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: bucketName,
		Object: destinationPath,
	}, src)
	if err != nil {
		return fmt.Errorf("error copying object: %v", err)
	}

	// Set permissions
	if permission == 1 { // Public
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": ["s3:GetObject"],
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Resource": ["arn:aws:s3:::%s/%s"],
					"Sid": ""
				}
			]
		}`, bucketName, destinationPath)

		err = client.SetObjectLockConfig(ctx, bucketName, &minio.ObjectLockConfig{
			ObjectLockEnabled: "Enabled",
		})
		if err != nil {
			return fmt.Errorf("error setting object lock config: %v", err)
		}
	}

	return nil
}

// RemoveFileOrDirectoryIfExists removes a file or directory if it exists
func (s *MinioStorageService) RemoveFileOrDirectoryIfExists(serverKey, bucketName, path string) error {
	client, exists := s.clients[serverKey]
	if !exists {
		return fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()

	// Check if the object exists
	exists, err := s.CheckFileOrDirectoryIsExists(serverKey, bucketName, path)
	if err != nil {
		return err
	}

	if !exists {
		return nil // Object doesn't exist, nothing to do
	}

	// Remove the object
	err = client.RemoveObject(ctx, bucketName, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("error removing object: %v", err)
	}

	return nil
}

// AddOrRemoveOrChangePermission changes the permission of a file or directory
func (s *MinioStorageService) AddOrRemoveOrChangePermission(serverKey, bucketName, path string, permission int) error {
	client, exists := s.clients[serverKey]
	if !exists {
		return fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()

	if permission == 1 { // Public
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": ["s3:GetObject"],
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Resource": ["arn:aws:s3:::%s/%s"],
					"Sid": ""
				}
			]
		}`, bucketName, path)

		err := client.SetObjectLockConfig(ctx, bucketName, &minio.ObjectLockConfig{
			ObjectLockEnabled: "Enabled",
		})
		if err != nil {
			return fmt.Errorf("error setting object lock config: %v", err)
		}
	} else {
		// For private, we remove any public access policies
		// This is simplified for demonstration. In a real implementation,
		// you would need to manage bucket policies more carefully.
	}

	return nil
}

// GenerateUrl generates a URL for a file
func (s *MinioStorageService) GenerateUrl(serverKey, bucketName, path string) string {
	client, exists := s.clients[serverKey]
	if !exists {
		return ""
	}

	// For public files, we can just construct the URL
	endpoint := s.config.Endpoints[serverKey]
	protocol := "http"
	if s.config.UseSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, bucketName, path)
}

// GeneratePreSignedUrl generates a pre-signed URL for a file with expiry
func (s *MinioStorageService) GeneratePreSignedUrl(serverKey, bucketName, path string, expiry time.Duration) (string, error) {
	client, exists := s.clients[serverKey]
	if !exists {
		return "", fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()

	// Generate presigned URL
	presignedURL, err := client.PresignedGetObject(ctx, bucketName, path, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("error generating presigned URL: %v", err)
	}

	return presignedURL.String(), nil
}

// DownloadFileOrDirectory downloads a file or directory
func (s *MinioStorageService) DownloadFileOrDirectory(serverKey, bucketName, path string) (io.Reader, error) {
	client, exists := s.clients[serverKey]
	if !exists {
		return nil, fmt.Errorf("server key %s not found", serverKey)
	}

	ctx := context.Background()

	// Get the object
	obj, err := client.GetObject(ctx, bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting object: %v", err)
	}

	return obj, nil
}

// Ensure MinioStorageService implements IStorageService
var _ contractStorage.IStorageService = &MinioStorageService{}
