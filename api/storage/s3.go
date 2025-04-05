package storage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	minioClient *minio.Client
}

type Action string

const (
	S3PUT Action = "PUT"
	S3GET Action = "GET"
)

func InitiateClient(endpoint, accessKeyID, secretAccessKey string) (*S3, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:       true, // Set to true if using HTTPS
		BucketLookup: minio.BucketLookupPath,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	s3 := new(S3)
	s3.minioClient = client

	return s3, nil
}

func (s3 S3) DownloadFile(bucket, objectKey, path, filename string) error {
	if s3.minioClient == nil {
		return fmt.Errorf("MinIO client not initiate yet")
	}

	err := os.MkdirAll(path+"/"+bucket, os.ModePerm)

	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create a file to write the downloaded content
	outFile, err := os.Create(path + "/" + bucket + "/" + filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	decoded_object_key, err := url.QueryUnescape(objectKey)

	if err != nil {
		return err
	}

	// Download the file
	err = s3.minioClient.FGetObject(context.Background(), bucket, decoded_object_key, path+"/"+bucket+"/"+filename, minio.GetObjectOptions{})

	if err != nil {
		os.Remove(path + "/" + bucket + "/" + filename)
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}

func (s3 S3) HeadFile(bucket, objectKey string) (*minio.ObjectInfo, error) {
	if s3.minioClient == nil {
		return nil, fmt.Errorf("minIO client not initiate yet")
	}

	if bucket == "" || objectKey == "" {
		return nil, fmt.Errorf("bucket / Object key cannot be empty")
	}

	decoded_object_key, err := url.QueryUnescape(objectKey)

	if err != nil {
		return nil, err
	}

	// Head the file
	objectInfo, err := s3.minioClient.StatObject(context.Background(), bucket, decoded_object_key, minio.StatObjectOptions{})

	return &objectInfo, err
}

func (s3 S3) GeneratePresignedUrl(action Action, bucket, objectKey string, expiration int64) (*url.URL, error) {
	if action == S3PUT {
		presignedURL, err := s3.minioClient.PresignedPutObject(context.Background(), bucket, objectKey, time.Duration(expiration)*time.Minute)

		return presignedURL, err
	} else if action == S3GET {
		presignedURL, err := s3.minioClient.PresignedGetObject(context.Background(), bucket, objectKey, time.Duration(expiration)*time.Minute, nil)

		return presignedURL, err
	}

	return nil, fmt.Errorf("invalid action: %s", action)
}
