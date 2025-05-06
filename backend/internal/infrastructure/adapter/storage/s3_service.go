package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type S3Handler struct {
	client *s3.Client
	bucket string
}

func NewS3Handler(cfg S3Config) (*S3Handler, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}
	return &S3Handler{
		client: s3.NewFromConfig(awsCfg),
		bucket: cfg.Bucket,
	}, nil
}

func (u *S3Handler) Upload(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &u.bucket,
		Key:    PointerOf(filepath.Base(path)),
		Body:   f,
		ACL:    types.ObjectCannedACLPrivate,
	})
	return err
}

// PointerOf returns a pointer to the provided value of any type(in this case the object path).
func PointerOf[T any](value T) *T {
	return &value
}

// StreamKeys streams object keys from S3
func (d *S3Handler) StreamKeys(ctx context.Context, bucket string) (<-chan string, error) {
	ch := make(chan string)

	go func() {
		defer close(ch)
		var continuationToken *string

		for {
			resp, err := d.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket:            &bucket,
				ContinuationToken: continuationToken,
			})
			if err != nil {
				return
			}

			for _, obj := range resp.Contents {
				ch <- *obj.Key
			}

			if resp.IsTruncated != nil && !*resp.IsTruncated {
				break
			}
			continuationToken = resp.NextContinuationToken
		}
	}()
	return ch, nil
}

// DownloadStream downloads and returns an io.ReadCloser from S3 for streaming
func (d *S3Handler) DownloadStream(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	resp, err := d.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
