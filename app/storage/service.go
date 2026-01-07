package storage

import "context"

type Service struct {
	store ObjectStore
}

func New(store ObjectStore) *Service {
	return &Service{store: store}
}

func (s *Service) StoreAvatar(ctx context.Context, userID string, data []byte) error {
	return s.store.Put(
		ctx,
		"avatars",
		userID+".png",
		data,
		"image/png",
	)
}
