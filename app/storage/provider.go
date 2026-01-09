package storage

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import "context"

// ------------------------------------------------------------------------------------------------
// Interfaces
// ------------------------------------------------------------------------------------------------

type ObjectStore interface {
	Put(ctx context.Context, bucket, key string, data []byte, contentType string) error
	Get(ctx context.Context, bucket, key string) ([]byte, error)
	Delete(ctx context.Context, bucket, key string) error
}

// ------------------------------------------------------------------------------------------------
