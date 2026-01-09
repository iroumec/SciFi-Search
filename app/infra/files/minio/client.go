package minio

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

func NewClient() (*s3.Client, error) {

	// Carga de configuración.
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("MINIO_ROOT_USER"),
				os.Getenv("MINIO_ROOT_PASSWORD"),
				"",
			),
		),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}

	// Devolución de cliente con la configuración definida.
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("MINIO_HOST"))
		o.UsePathStyle = true
	}), nil
}

// ------------------------------------------------------------------------------------------------
