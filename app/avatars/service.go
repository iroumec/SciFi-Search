package avatars

// ------------------------------------------------------------------

import (
	"context"
	"fmt"
)

// ------------------------------------------------------------------

type Storage interface {
	Put(ctx context.Context, bucket, key string, data []byte, contentType string) error
	Delete(ctx context.Context, bucket, key string) error
}

// ------------------------------------------------------------------

type Service struct {
	store     Storage
	bucket    string
	publicURL string
}

// ------------------------------------------------------------------

func New(store Storage, bucket, publicURL string) *Service {
	return &Service{
		store:     store,
		bucket:    bucket,
		publicURL: publicURL,
	}
}

// ------------------------------------------------------------------

func (s *Service) Upload(ctx context.Context, userID int32, data []byte) (string, error) {
	key := fmt.Sprintf("%d.jpg", userID)

	if err := s.store.Put(ctx, s.bucket, key, data, "image/jpeg"); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, key), nil
}

// ------------------------------------------------------------------

func (s *Service) Delete(ctx context.Context, userID int32) error {
	key := fmt.Sprintf("%d.jpg", userID)
	return s.store.Delete(ctx, s.bucket, key)
}

// ------------------------------------------------------------------
