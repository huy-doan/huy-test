package storage

import (
	"context"
	"io"
)

type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type S3Service interface {
	Upload(ctx context.Context, path string) error

	UploadStream(ctx context.Context, key string, body io.Reader) error

	UploadStreamWithContentLength(ctx context.Context, key string, body io.Reader, contentLength int64) error

	StreamKeys(ctx context.Context, bucket string) (<-chan string, error)

	DownloadStream(ctx context.Context, bucket, key string) (io.ReadCloser, error)

	GetS3KeyFromRemotePath(remotePath, remoteDir string) string

	GetObjectImportStatus(ctx context.Context, key string) (string, error)

	SetObjectImportStatus(ctx context.Context, key, status string) error
}
