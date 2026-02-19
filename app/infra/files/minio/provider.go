package minio

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"scifi-search/app/storage"
)

// ------------------------------------------------------------------------------------------------
// Structures
// ------------------------------------------------------------------------------------------------

type Provider struct {
	client *s3.Client
}

// ------------------------------------------------------------------------------------------------

func (p *Provider) Put(
	ctx context.Context,
	bucket, key string,
	data []byte,
	contentType string,
) error {

	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})

	return err
}

// ------------------------------------------------------------------------------------------------

func (p *Provider) Get(ctx context.Context, bucket, key string) ([]byte, error) {
	out, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

// ------------------------------------------------------------------------------------------------

func (p *Provider) Delete(ctx context.Context, bucket, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	return err
}

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Takes a client and returns a provider.
func New(client *s3.Client) storage.ObjectStore {
	return &Provider{client: client}
}

// ------------------------------------------------------------------------------------------------
