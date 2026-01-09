package bootstrap

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"context"
	"log"
	"os"
	"scifi-search/app/avatars"
	"scifi-search/app/infra/files/minio"
)

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func getAvatarsService() *avatars.Service {

	minioClient, err := minio.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// SetUp de Bucket.
	if err := minio.EnsurePublicBucket(
		context.Background(),
		minioClient,
		"avatars",
	); err != nil {
		log.Fatal(err)
	}

	// Storage provider.
	objectStore := minio.New(minioClient)

	// Dominio.
	avatarService := avatars.New(
		objectStore,
		"avatars",
		os.Getenv("MINIO_PUBLIC_URL"),
	)

	return avatarService
}

// ------------------------------------------------------------------------------------------------
