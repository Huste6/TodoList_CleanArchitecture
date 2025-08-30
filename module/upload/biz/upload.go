package uploadbiz

import (
	"context"
	"fmt"
	"g09/common"
	"path/filepath"
	"strings"
	"time"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) (*common.Image, error)
	DeleteImage(ctx context.Context, publicID string) error
	ExtractPublicID(imageURL string) string
}

type uploadBiz struct {
	provider UploadProvider
}

func NewUploadBiz(provider UploadProvider) *uploadBiz {
	return &uploadBiz{
		provider: provider,
	}
}

func (biz *uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string) (*common.Image, error) {
	fileExt := filepath.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, fileExt)

	dst := fmt.Sprintf("%s/%d_%s%s", folder, time.Now().UTC().Unix(), fileName, fileExt)

	img, err := biz.provider.SaveFileUploaded(ctx, data, dst)
	if err != nil {
		return nil, common.ErrCannotCreateEntity("image", err)
	}

	return img, nil
}

func (biz *uploadBiz) DeleteImage(ctx context.Context, imageURL string) error {
	if imageURL == "" {
		return nil
	}

	publicID := biz.provider.ExtractPublicID(imageURL)
	if publicID == "" {
		return common.ErrInvalidRequest(fmt.Errorf("invalid image URL"))
	}

	if err := biz.provider.DeleteImage(ctx, publicID); err != nil {
		return common.ErrCannotDeleteEntity("image", err)
	}

	return nil
}
