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
	"github.com/huydq/test/internal/pkg/utils"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func NewS3Client(cfg S3Config) (S3Service, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}
	return &S3Client{
		client: s3.NewFromConfig(awsCfg),
		bucket: cfg.Bucket,
	}, nil
}

func (u *S3Client) Upload(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &u.bucket,
		Key:    utils.ToPtr(filepath.Base(path)),
		Body:   f,
		ACL:    types.ObjectCannedACLPrivate,
	})
	return err
}

func (u *S3Client) UploadStream(ctx context.Context, key string, body io.Reader) error {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &u.bucket,
		Key:    utils.ToPtr(key),
		Body:   body,
		ACL:    types.ObjectCannedACLPrivate,
	})
	return err
}

func (u *S3Client) UploadStreamWithContentLength(ctx context.Context, key string, body io.Reader, contentLength int64) error {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &u.bucket,
		Key:           utils.ToPtr(key),
		Body:          body,
		ACL:           types.ObjectCannedACLPrivate,
		ContentLength: &contentLength,
	})
	return err
}

// StreamKeys streams object keys from S3
func (d *S3Client) StreamKeys(ctx context.Context, bucket string) (<-chan string, error) {
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
func (d *S3Client) DownloadStream(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	resp, err := d.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// GetS3KeyFromRemotePath generates an S3 key by removing the remote directory prefix from the remote path.
func (d *S3Client) GetS3KeyFromRemotePath(remotePath, remoteDir string) string {
	key := remotePath
	if len(key) > 0 && key[0] == '/' {
		key = key[1:]
	}
	if remoteDir != "" {
		prefix := remoteDir
		if len(prefix) > 0 && prefix[0] == '/' {
			prefix = prefix[1:]
		}
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			key = key[len(prefix):]
			if len(key) > 0 && key[0] == '/' {
				key = key[1:]
			}
		}
	}
	return key
}

// GetObjectImportStatus returns the value of the 'import_status' tag for the given S3 key, or empty string if not set.
func (d *S3Client) GetObjectImportStatus(ctx context.Context, key string) (string, error) {
	resp, err := d.client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
		Bucket: &d.bucket,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}
	for _, tag := range resp.TagSet {
		if tag.Key != nil && *tag.Key == "import_status" && tag.Value != nil {
			return *tag.Value, nil
		}
	}
	return "", nil
}

// SetObjectImportStatus sets the 'import_status' tag for the given S3 key.
func (d *S3Client) SetObjectImportStatus(ctx context.Context, key, status string) error {
	_, err := d.client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
		Bucket: &d.bucket,
		Key:    &key,
		Tagging: &types.Tagging{
			TagSet: []types.Tag{{
				Key:   utils.ToPtr("import_status"),
				Value: utils.ToPtr(status),
			}},
		},
	})
	return err
}
