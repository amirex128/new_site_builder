package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/amirex128/new_site_builder/internal/domain/enums"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageClient represents an S3-compatible storage client
type StorageClient struct {
	client *minio.Client
	host   string
}

// NewStorageClient creates a new storage client
func NewStorageClient(host, accessKey, secretKey string) *StorageClient {
	// Initialize MinIO client
	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize storage client: %v", err))
	}

	return &StorageClient{
		client: client,
		host:   host,
	}
}

// StorageService implements IStorageService
type StorageService struct {
	clients map[string]*StorageClient
}

// NewStorageService creates a new instance of StorageService with multiple clients
func NewStorageService(s1Client, s2Client, s3Client *StorageClient) *StorageService {
	clients := make(map[string]*StorageClient)
	clients["S1"] = s1Client
	clients["S2"] = s2Client
	clients["S3"] = s3Client

	return &StorageService{
		clients: clients,
	}
}

// getClient returns the client for the given key
func (s *StorageService) getClient(key string) (*StorageClient, error) {
	client, ok := s.clients[key]
	if !ok {
		return nil, fmt.Errorf("client with key %s not found", key)
	}
	return client, nil
}

// CreateFileOrDirectoryIfNotExists creates a file or directory if it doesn't exist
func (s *StorageService) CreateFileOrDirectoryIfNotExists(serverKey, bucketName, key string, permission enums.FileItemPermissionEnum, fileStream ...io.Reader) (string, error) {
	client, err := s.getClient(serverKey)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	var reader io.Reader
	var size int64

	if len(fileStream) > 0 && fileStream[0] != nil {
		reader = fileStream[0]

		// Try to determine content size if reader supports it
		if seeker, ok := reader.(io.Seeker); ok {
			currentPos, err := seeker.Seek(0, io.SeekCurrent)
			if err != nil {
				return "", err
			}
			endPos, err := seeker.Seek(0, io.SeekEnd)
			if err != nil {
				return "", err
			}
			_, err = seeker.Seek(currentPos, io.SeekStart)
			if err != nil {
				return "", err
			}
			size = endPos - currentPos
		}
	} else {
		// For directories, create an empty object with trailing slash
		if !strings.HasSuffix(key, "/") {
			key += "/"
		}
		reader = bytes.NewReader([]byte{})
		size = 0
	}

	_, err = client.client.PutObject(ctx, bucketName, key, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	// Set permissions
	err = s.AddOrRemoveOrChangePermission(serverKey, bucketName, key, permission, false)
	if err != nil {
		return "", err
	}

	return s.GenerateUrl(serverKey, bucketName, key), nil
}

// CheckFileOrDirectoryIsExists checks if a file or directory exists
func (s *StorageService) CheckFileOrDirectoryIsExists(serverKey, bucketName, key string) (bool, error) {
	client, err := s.getClient(serverKey)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	_, err = client.client.StatObject(ctx, bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		// Check if error is "not found"
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// RemoveFileOrDirectoryIfExists removes a file or directory if it exists
func (s *StorageService) RemoveFileOrDirectoryIfExists(serverKey, bucketName, key string) (bool, error) {
	client, err := s.getClient(serverKey)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	exists, err := s.CheckFileOrDirectoryIsExists(serverKey, bucketName, key)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	err = client.client.RemoveObject(ctx, bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateBucketIfNotExists creates a bucket if it doesn't exist
func (s *StorageService) CreateBucketIfNotExists(serverKey, bucketName string) error {
	client, err := s.getClient(serverKey)
	if err != nil {
		return err
	}

	ctx := context.Background()
	exists, err := client.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return client.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}

// DeleteBucketIfExists deletes a bucket if it exists
func (s *StorageService) DeleteBucketIfExists(serverKey, bucketName string) error {
	client, err := s.getClient(serverKey)
	if err != nil {
		return err
	}

	ctx := context.Background()
	exists, err := client.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if exists {
		return client.client.RemoveBucket(ctx, bucketName)
	}

	return nil
}

// DownloadFileOrDirectory downloads a file or directory
func (s *StorageService) DownloadFileOrDirectory(serverKey, bucketName, key string) (io.ReadCloser, error) {
	client, err := s.getClient(serverKey)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	obj, err := client.client.GetObject(ctx, bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// GeneratePreSignedUrl generates a pre-signed URL for a file or directory
func (s *StorageService) GeneratePreSignedUrl(serverKey, bucketName, key string, expiry time.Duration) (string, error) {
	client, err := s.getClient(serverKey)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	presignedURL, err := client.client.PresignedGetObject(ctx, bucketName, key, expiry, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// GenerateUrl generates a URL for a file or directory
func (s *StorageService) GenerateUrl(serverKey, bucketName, key string) string {
	client, err := s.getClient(serverKey)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("https://%s/%s/%s", client.host, bucketName, key)
}

// RenameOrMoveFileOrDirectory renames or moves a file or directory
func (s *StorageService) RenameOrMoveFileOrDirectory(serverKey, bucketName, oldKey, newKey string, currentPolicy enums.FileItemPermissionEnum) (string, error) {
	client, err := s.getClient(serverKey)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	// Check if source exists
	exists, err := s.CheckFileOrDirectoryIsExists(serverKey, bucketName, oldKey)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("old file or directory does not exist")
	}

	// Check if destination already exists
	exists, err = s.CheckFileOrDirectoryIsExists(serverKey, bucketName, newKey)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("new file or directory already exists")
	}

	// If it's a directory (ends with '/')
	if strings.HasSuffix(oldKey, "/") {
		// List all objects with the prefix
		objectCh := client.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Prefix:    oldKey,
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				return "", object.Err
			}

			// Create the new key by replacing the old prefix with the new one
			destKey := strings.Replace(object.Key, oldKey, newKey, 1)

			// Copy the object
			_, err := client.client.CopyObject(ctx, minio.CopyDestOptions{
				Bucket: bucketName,
				Object: destKey,
			}, minio.CopySrcOptions{
				Bucket: bucketName,
				Object: object.Key,
			})
			if err != nil {
				return "", err
			}

			// Update permissions for the new object
			err = s.AddOrRemoveOrChangePermission(serverKey, bucketName, destKey, currentPolicy, false)
			if err != nil {
				return "", err
			}

			// Remove the old object
			err = client.client.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				return "", err
			}
		}
	} else {
		// Copy the object
		_, err := client.client.CopyObject(ctx, minio.CopyDestOptions{
			Bucket: bucketName,
			Object: newKey,
		}, minio.CopySrcOptions{
			Bucket: bucketName,
			Object: oldKey,
		})
		if err != nil {
			return "", err
		}

		// Update permissions
		err = s.AddOrRemoveOrChangePermission(serverKey, bucketName, oldKey, currentPolicy, true)
		if err != nil {
			return "", err
		}

		err = s.AddOrRemoveOrChangePermission(serverKey, bucketName, newKey, currentPolicy, false)
		if err != nil {
			return "", err
		}

		// Remove the old object
		err = client.client.RemoveObject(ctx, bucketName, oldKey, minio.RemoveObjectOptions{})
		if err != nil {
			return "", err
		}
	}

	return s.GenerateUrl(serverKey, bucketName, newKey), nil
}

// CopyFileOrDirectory copies a file or directory
func (s *StorageService) CopyFileOrDirectory(serverKey, bucketName, sourceKey, destinationDirectory string, currentPolicy enums.FileItemPermissionEnum) (bool, error) {
	client, err := s.getClient(serverKey)
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	// Check if source exists
	exists, err := s.CheckFileOrDirectoryIsExists(serverKey, bucketName, sourceKey)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("source file or directory does not exist")
	}

	// Check if destination already exists
	exists, err = s.CheckFileOrDirectoryIsExists(serverKey, bucketName, destinationDirectory)
	if err != nil {
		return false, err
	}
	if exists {
		return false, fmt.Errorf("destination directory already exists")
	}

	// Construct the new key
	fileName := filepath.Base(sourceKey)
	newKey := strings.TrimSuffix(destinationDirectory, "/") + "/" + fileName

	// If it's a directory (ends with '/')
	if strings.HasSuffix(sourceKey, "/") {
		// List all objects with the prefix
		objectCh := client.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Prefix:    sourceKey,
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				return false, object.Err
			}

			// Create the new key by replacing the old prefix with the new one
			destKey := strings.Replace(object.Key, sourceKey, newKey, 1)

			// Copy the object
			_, err := client.client.CopyObject(ctx, minio.CopyDestOptions{
				Bucket: bucketName,
				Object: destKey,
			}, minio.CopySrcOptions{
				Bucket: bucketName,
				Object: object.Key,
			})
			if err != nil {
				return false, err
			}

			// Update permissions for the new object
			err = s.AddOrRemoveOrChangePermission(serverKey, bucketName, destKey, currentPolicy, false)
			if err != nil {
				return false, err
			}
		}
	} else {
		// Copy the object
		_, err := client.client.CopyObject(ctx, minio.CopyDestOptions{
			Bucket: bucketName,
			Object: newKey,
		}, minio.CopySrcOptions{
			Bucket: bucketName,
			Object: sourceKey,
		})
		if err != nil {
			return false, err
		}

		// Update permissions
		err = s.AddOrRemoveOrChangePermission(serverKey, bucketName, newKey, currentPolicy, false)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// AddOrRemoveOrChangePermission adds, removes, or changes permissions on a file or directory
func (s *StorageService) AddOrRemoveOrChangePermission(serverKey, bucketName, key string, permission enums.FileItemPermissionEnum, justRemove bool) error {
	client, err := s.getClient(serverKey)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Resource ARN
	resourceArn := fmt.Sprintf("arn:aws:s3:::%s/%s", bucketName, key)

	// Get current bucket policy
	var policyJSON string
	policyJSON, err = client.client.GetBucketPolicy(ctx, bucketName)
	if err != nil {
		// If error is "policy not found", create a new one
		if minio.ToErrorResponse(err).Code == "NoSuchBucketPolicy" {
			policyJSON = `{"Version":"2012-10-17","Statement":[]}`
		} else {
			return err
		}
	}

	// Parse policy JSON
	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		return err
	}

	// Initialize statements if not exists
	if _, ok := policy["Statement"]; !ok {
		policy["Statement"] = []interface{}{}
	}

	statements := policy["Statement"].([]interface{})

	// Find public and private statements
	var publicStatement map[string]interface{}
	var privateStatement map[string]interface{}
	var publicIndex, privateIndex int
	var publicResources, privateResources []interface{}

	for i, stmt := range statements {
		statement := stmt.(map[string]interface{})
		if sid, ok := statement["Sid"].(string); ok {
			if sid == "PublicRead" {
				publicStatement = statement
				publicIndex = i
				if resources, ok := statement["Resource"].([]interface{}); ok {
					publicResources = resources
				} else {
					publicResources = []interface{}{}
				}
			} else if sid == "DenyAllExceptOwner" {
				privateStatement = statement
				privateIndex = i
				if resources, ok := statement["Resource"].([]interface{}); ok {
					privateResources = resources
				} else {
					privateResources = []interface{}{}
				}
			}
		}
	}

	// If just removing, remove from both lists
	if justRemove {
		// Remove from public resources
		if publicStatement != nil {
			newPublicResources := []interface{}{}
			for _, res := range publicResources {
				if res.(string) != resourceArn {
					newPublicResources = append(newPublicResources, res)
				}
			}
			publicResources = newPublicResources
			publicStatement["Resource"] = publicResources

			// If no resources left, remove the statement
			if len(publicResources) == 0 {
				statements = append(statements[:publicIndex], statements[publicIndex+1:]...)
				// Adjust privateIndex if needed
				if privateIndex > publicIndex {
					privateIndex--
				}
			}
		}

		// Remove from private resources
		if privateStatement != nil {
			newPrivateResources := []interface{}{}
			for _, res := range privateResources {
				if res.(string) != resourceArn {
					newPrivateResources = append(newPrivateResources, res)
				}
			}
			privateResources = newPrivateResources
			privateStatement["Resource"] = privateResources

			// If no resources left, remove the statement
			if len(privateResources) == 0 {
				if privateIndex < len(statements) {
					statements = append(statements[:privateIndex], statements[privateIndex+1:]...)
				}
			}
		}
	} else {
		// Add to appropriate list based on permission
		if permission == enums.FileItemPublicPermission {
			// Add to public, remove from private

			// Create public statement if it doesn't exist
			if publicStatement == nil {
				publicStatement = map[string]interface{}{
					"Sid":       "PublicRead",
					"Effect":    "Allow",
					"Principal": map[string]string{"AWS": "*"},
					"Action":    []string{"s3:GetObject"},
					"Resource":  []interface{}{},
				}
				publicResources = []interface{}{}
				statements = append(statements, publicStatement)
			}

			// Add to public resources if not already there
			found := false
			for _, res := range publicResources {
				if res.(string) == resourceArn {
					found = true
					break
				}
			}
			if !found {
				publicResources = append(publicResources, resourceArn)
			}
			publicStatement["Resource"] = publicResources

			// Remove from private resources
			if privateStatement != nil {
				newPrivateResources := []interface{}{}
				for _, res := range privateResources {
					if res.(string) != resourceArn {
						newPrivateResources = append(newPrivateResources, res)
					}
				}
				privateResources = newPrivateResources
				privateStatement["Resource"] = privateResources

				// If no resources left, remove the statement
				if len(privateResources) == 0 {
					for i, stmt := range statements {
						if stmt.(map[string]interface{})["Sid"].(string) == "DenyAllExceptOwner" {
							statements = append(statements[:i], statements[i+1:]...)
							break
						}
					}
				}
			}
		} else {
			// Add to private, remove from public

			// Create private statement if it doesn't exist
			if privateStatement == nil {
				privateStatement = map[string]interface{}{
					"Sid":       "DenyAllExceptOwner",
					"Effect":    "Deny",
					"Principal": map[string]string{"AWS": "*"},
					"Action":    []string{"s3:GetObject"},
					"Resource":  []interface{}{},
				}
				privateResources = []interface{}{}
				statements = append(statements, privateStatement)
			}

			// Add to private resources if not already there
			found := false
			for _, res := range privateResources {
				if res.(string) == resourceArn {
					found = true
					break
				}
			}
			if !found {
				privateResources = append(privateResources, resourceArn)
			}
			privateStatement["Resource"] = privateResources

			// Remove from public resources
			if publicStatement != nil {
				newPublicResources := []interface{}{}
				for _, res := range publicResources {
					if res.(string) != resourceArn {
						newPublicResources = append(newPublicResources, res)
					}
				}
				publicResources = newPublicResources
				publicStatement["Resource"] = publicResources

				// If no resources left, remove the statement
				if len(publicResources) == 0 {
					for i, stmt := range statements {
						if stmt.(map[string]interface{})["Sid"].(string) == "PublicRead" {
							statements = append(statements[:i], statements[i+1:]...)
							break
						}
					}
				}
			}
		}
	}

	// Update policy with modified statements
	policy["Statement"] = statements

	// Convert policy back to JSON
	updatedPolicyJSON, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	// Set the updated bucket policy
	return client.client.SetBucketPolicy(ctx, bucketName, string(updatedPolicyJSON))
}
