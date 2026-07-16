package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
)

type MinIOStorage struct {
	logger *slog.Logger
	Client *minio.Client
	bucket string
}

func NewMinIOStorage(l *slog.Logger, c *minio.Client, b string, ctx context.Context) (*MinIOStorage, error) {
	exists, err := c.BucketExists(ctx, b)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = c.MakeBucket(
			ctx,
			b,
			minio.MakeBucketOptions{},
		)

		if err != nil {
			return nil, err
		}

		// Make bucket publicly readable
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"AWS": ["*"]
					},
					"Action": [
						"s3:GetObject"
					],
					"Resource": [
						"arn:aws:s3:::%s/*"
					]
				}
			]
		}`, b)

		err = c.SetBucketPolicy(
			ctx,
			b,
			policy,
		)

		if err != nil {
			return nil, err
		}
	}

	return &MinIOStorage{logger: l, Client: c, bucket: b}, nil
}

func (m *MinIOStorage) Save(ctx context.Context, key string, r io.Reader) error {
	_, err := m.Client.PutObject(
		ctx,
		m.bucket,
		key,
		r,
		-1,
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		},
	)

	return err
}

func (m *MinIOStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := m.Client.GetObject(
		ctx,
		m.bucket,
		key,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, err
	}

	return object, nil
}

func (m *MinIOStorage) Delete(ctx context.Context, key string) error {
	return m.Client.RemoveObject(
		ctx,
		m.bucket,
		key,
		minio.RemoveObjectOptions{},
	)
}

//Save(ctx context.Context, key string, r io.Reader) error
// Get(ctx context.Context, key string) (io.ReadCloser, error)
// Delete(ctx context.Context, key string) error
