package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

var ErrBlobNotFound = errors.New("blob not found")

type BlobStore interface {
	Put(ctx context.Context, key string, r io.Reader, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
}

func validateBlobKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("blob key cannot be empty")
	}

	if strings.HasPrefix(key, "/") || filepath.IsAbs(key) || strings.Contains(key, "\\") {
		return "", fmt.Errorf("invalid blob key %q", key)
	}

	for _, part := range strings.Split(key, "/") {
		if part == ".." {
			return "", fmt.Errorf("invalid blob key %q", key)
		}
	}

	cleaned := path.Clean(key)
	if cleaned == "." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf("invalid blob key %q", key)
	}

	return cleaned, nil
}

type LocalStore struct {
	root string
}

func NewLocalStore(root string) (*LocalStore, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local blob path: %w", err)
	}

	return &LocalStore{root: absolute}, nil
}

func (s *LocalStore) resolve(key string) (string, error) {
	cleaned, err := validateBlobKey(key)
	if err != nil {
		return "", err
	}

	return filepath.Join(s.root, filepath.FromSlash(cleaned)), nil
}

func (s *LocalStore) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
	target, err := s.resolve(key)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("failed to create blob directory: %w", err)
	}

	temp, err := os.CreateTemp(filepath.Dir(target), ".blob-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary blob: %w", err)
	}
	temp_name := temp.Name()

	if _, err := io.Copy(temp, r); err != nil {
		temp.Close()
		os.Remove(temp_name)
		return err
	}

	if err := temp.Close(); err != nil {
		os.Remove(temp_name)
		return err
	}

	if err := os.Rename(temp_name, target); err != nil {
		os.Remove(temp_name)
		return fmt.Errorf("failed to store blob: %w", err)
	}

	return nil
}

func (s *LocalStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	target, err := s.resolve(key)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(target)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrBlobNotFound
	} else if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *LocalStore) Exists(ctx context.Context, key string) (bool, error) {
	target, err := s.resolve(key)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(target)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (s *LocalStore) Delete(ctx context.Context, key string) error {
	target, err := s.resolve(key)
	if err != nil {
		return err
	}

	if err := os.Remove(target); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	return nil
}

func (s *LocalStore) DeletePrefix(ctx context.Context, prefix string) error {
	cleaned, err := validateBlobKey(prefix)
	if err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(s.root, filepath.FromSlash(cleaned)))
}

type S3BlobStore struct {
	client *minio.Client
	bucket string
}

func (s *S3BlobStore) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
	cleaned, err := validateBlobKey(key)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(ctx, s.bucket, cleaned, r, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

func (s *S3BlobStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	cleaned, err := validateBlobKey(key)
	if err != nil {
		return nil, err
	}

	if _, err := s.client.StatObject(ctx, s.bucket, cleaned, minio.StatObjectOptions{}); err != nil {
		if isMissingObject(err) {
			return nil, ErrBlobNotFound
		}
		return nil, err
	}

	object, err := s.client.GetObject(ctx, s.bucket, cleaned, minio.GetObjectOptions{})
	if err != nil {
		if isMissingObject(err) {
			return nil, ErrBlobNotFound
		}
		return nil, err
	}

	return object, nil
}

func (s *S3BlobStore) Exists(ctx context.Context, key string) (bool, error) {
	cleaned, err := validateBlobKey(key)
	if err != nil {
		return false, err
	}

	_, err = s.client.StatObject(ctx, s.bucket, cleaned, minio.StatObjectOptions{})
	if isMissingObject(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (s *S3BlobStore) Delete(ctx context.Context, key string) error {
	cleaned, err := validateBlobKey(key)
	if err != nil {
		return err
	}

	return s.client.RemoveObject(ctx, s.bucket, cleaned, minio.RemoveObjectOptions{})
}

func (s *S3BlobStore) DeletePrefix(ctx context.Context, prefix string) error {
	cleaned, err := validateBlobKey(prefix)
	if err != nil {
		return err
	}

	objects := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    cleaned,
		Recursive: true,
	})

	for object := range objects {
		if object.Err != nil {
			return object.Err
		}

		if err := s.client.RemoveObject(ctx, s.bucket, object.Key, minio.RemoveObjectOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func isMissingObject(err error) bool {
	if err == nil {
		return false
	}

	response := minio.ToErrorResponse(err)
	return response.Code == "NoSuchKey" || response.Code == "NotFound" || response.StatusCode == 404
}

func NewBlobStore(backend, localPath string, s3 *S3, bucket string) (BlobStore, error) {
	switch strings.ToLower(strings.TrimSpace(backend)) {
	case "", "local":
		return NewLocalStore(localPath)
	case "s3":
		if s3 == nil || s3.minioClient == nil {
			return nil, fmt.Errorf("S3 blob backend requires an initialized S3 client")
		}
		if bucket == "" {
			return nil, fmt.Errorf("S3 blob backend requires S3_BUCKET")
		}
		return &S3BlobStore{client: s3.minioClient, bucket: bucket}, nil
	default:
		return nil, fmt.Errorf("unsupported BLOB_BACKEND %q", backend)
	}
}
