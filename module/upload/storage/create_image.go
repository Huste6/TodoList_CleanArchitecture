package uploadstorage

import (
	"context"
	"g09/common"
)

func (s *uploadStore) CreateImage(ctx context.Context, data *common.Image) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
