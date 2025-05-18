package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageService implements IStorageService
type StorageService struct {
	client     *minio.Client
	bucketName string
}

// NewStorageService creates a new instance of StorageService
func NewStorageService(bucketName, region, accessKey, secretKey string) *StorageService {
	// Initialize MinIO client
	endpoint := "s3.amazonaws.com" // Default S3 endpoint
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
		Region: region,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize storage client: %v", err))
	}

	return &StorageService{
		client:     client,
		bucketName: bucketName,
	}
}

// CreateBucketIfNotExists creates a bucket if it doesn't exist
func (s *StorageService) CreateBucketIfNotExists(serverKey, bucketName string) error {
	ctx := context.Background()
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: serverKey,
		})
	}
	return nil
}

// CreateFileOrDirectoryIfNotExists creates a file or directory if it doesn't exist
func (s *StorageService) CreateFileOrDirectoryIfNotExists(serverKey, bucketName, path string, permission int, content ...io.Reader) error {
	ctx := context.Background()
	exists, err := s.CheckFileOrDirectoryIsExists(serverKey, bucketName, path)
	if err != nil {
		return err
	}
	if !exists {
		var reader io.Reader
		if len(content) > 0 && content[0] != nil {
			reader = content[0]
		} else {
			reader = bytes.NewReader([]byte{})
		}

		// Determine content size
		var size int64
		if seeker, ok := reader.(io.Seeker); ok {
			currentPos, err := seeker.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
			endPos, err := seeker.Seek(0, io.SeekEnd)
			if err != nil {
				return err
			}
			_, err = seeker.Seek(currentPos, io.SeekStart)
			if err != nil {
				return err
			}
			size = endPos - currentPos
		}

		_, err = s.client.PutObject(ctx, bucketName, path, reader, size, minio.PutObjectOptions{})
		return err
	}
	return nil
}

// CheckFileOrDirectoryIsExists checks if a file or directory exists
func (s *StorageService) CheckFileOrDirectoryIsExists(serverKey, bucketName, path string) (bool, error) {
	ctx := context.Background()
	_, err := s.client.StatObject(ctx, bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// RenameOrMoveFileOrDirectory renames or moves a file or directory
func (s *StorageService) RenameOrMoveFileOrDirectory(serverKey, bucketName, sourcePath, destinationPath string, permission int) error {
	ctx := context.Background()

	// Copy the object to the new location
	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: sourcePath,
	}
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: destinationPath,
	}
	_, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return err
	}

	// Delete the original object
	return s.client.RemoveObject(ctx, bucketName, sourcePath, minio.RemoveObjectOptions{})
}

// CopyFileOrDirectory copies a file or directory
func (s *StorageService) CopyFileOrDirectory(serverKey, bucketName, sourcePath, destinationPath string, permission int) error {
	ctx := context.Background()

	// Copy the object to the new location
	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: sourcePath,
	}
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: destinationPath,
	}
	_, err := s.client.CopyObject(ctx, dst, src)
	return err
}

// RemoveFileOrDirectoryIfExists removes a file or directory if it exists
func (s *StorageService) RemoveFileOrDirectoryIfExists(serverKey, bucketName, path string) error {
	ctx := context.Background()
	exists, err := s.CheckFileOrDirectoryIsExists(serverKey, bucketName, path)
	if err != nil {
		return err
	}
	if exists {
		return s.client.RemoveObject(ctx, bucketName, path, minio.RemoveObjectOptions{})
	}
	return nil
}

// AddOrRemoveOrChangePermission adds, removes, or changes permissions on a file or directory
func (s *StorageService) AddOrRemoveOrChangePermission(serverKey, bucketName, path string, permission int) error {
	// MinIO doesn't directly support Unix-style permissions
	// This is a placeholder implementation
	return nil
}

// GenerateUrl generates a URL for a file or directory
func (s *StorageService) GenerateUrl(serverKey, bucketName, path string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, path)
}

// GeneratePreSignedUrl generates a pre-signed URL for a file or directory
func (s *StorageService) GeneratePreSignedUrl(serverKey, bucketName, path string, expiry time.Duration) (string, error) {
	ctx := context.Background()
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, bucketName, path, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// DownloadFileOrDirectory downloads a file or directory
func (s *StorageService) DownloadFileOrDirectory(serverKey, bucketName, path string) (io.Reader, error) {
	ctx := context.Background()
	object, err := s.client.GetObject(ctx, bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Upload uploads a file to storage (convenience method)
func (s *StorageService) Upload(ctx context.Context, path string, content []byte, contentType string) (string, error) {
	reader := bytes.NewReader(content)
	_, err := s.client.PutObject(ctx, s.bucketName, path, reader, int64(len(content)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Generate URL for the uploaded file
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, path)
	return fileURL, nil
}

// Delete deletes a file from storage (convenience method)
func (s *StorageService) Delete(ctx context.Context, path string) error {
	return s.client.RemoveObject(ctx, s.bucketName, path, minio.RemoveObjectOptions{})
}

// GetSignedURL generates a pre-signed URL for the given path (convenience method)
func (s *StorageService) GetSignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	// Get presigned URL for object download
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, path, expiry, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// Download downloads a file from storage (convenience method)
func (s *StorageService) Download(ctx context.Context, path string) ([]byte, error) {
	// Get object
	object, err := s.client.GetObject(ctx, s.bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	// Read the object content
	content, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return content, nil
}
