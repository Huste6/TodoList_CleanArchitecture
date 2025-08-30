package uploadstorage

import (
	"context"
	"g09/common"
)

func (s *uploadStore) GetImage(ctx context.Context, id int) (*common.Image, error) {
	var image common.Image
	if err := s.db.Where("id = ?", id).First(&image).Error; err != nil {
		return nil, common.ErrCannotGetEntity("image", err)
	}
	return &image, nil
}
