package uploadprovider

import (
	"bytes"
	"context"
	"g09/common"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type cloudinaryProvider struct {
	cld    *cloudinary.Cloudinary
	folder string
}

func NewCloudinaryProvider(cloudName, apiKey, apiSecret, uploadFolder string) *cloudinaryProvider {
	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	return &cloudinaryProvider{
		cld:    cld,
		folder: uploadFolder,
	}
}

func (provider *cloudinaryProvider) SaveFileUploaded(ctx context.Context, data []byte, dst string) (*common.Image, error) {
	uploadResult, err := provider.cld.Upload.Upload(ctx, bytes.NewReader(data), uploader.UploadParams{
		PublicID:       dst,
		Folder:         provider.folder,
		ResourceType:   "image",
		AllowedFormats: []string{"jpg", "jpeg", "png", "gif", "webp", "bmp", "tiff"},
	})

	if err != nil {
		return nil, err
	}

	img := &common.Image{
		Url:       uploadResult.SecureURL,
		Width:     uploadResult.Width,
		Height:    uploadResult.Height,
		CloudName: provider.cld.Config.Cloud.CloudName,
		Extension: uploadResult.Format,
	}

	return img, nil
}

// DeleteImage deletes image from Cloudinary
func (provider *cloudinaryProvider) DeleteImage(ctx context.Context, publicID string) error {
	_, err := provider.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})
	return err
}

// ExtractPublicID extracts public ID from Cloudinary URL
func (provider *cloudinaryProvider) ExtractPublicID(imageURL string) string {
	// URL format: https://res.cloudinary.com/cloud_name/image/upload/v123456/folder/public_id.ext
	// hoặc: https://res.cloudinary.com/cloud_name/image/upload/folder/public_id.ext

	parts := strings.Split(imageURL, "/")
	if len(parts) < 7 {
		return ""
	}

	// Tìm vị trí "upload"
	uploadIndex := -1
	for i, part := range parts {
		if part == "upload" {
			uploadIndex = i
			break
		}
	}

	if uploadIndex == -1 || uploadIndex >= len(parts)-1 {
		return ""
	}

	// Lấy phần sau "upload", bỏ qua version nếu có (vXXXXXX)
	remainingParts := parts[uploadIndex+1:]

	// Bỏ qua version nếu part đầu tiên bắt đầu với "v" và là số
	if len(remainingParts) > 0 && strings.HasPrefix(remainingParts[0], "v") {
		if len(remainingParts[0]) > 1 {
			remainingParts = remainingParts[1:]
		}
	}

	if len(remainingParts) == 0 {
		return ""
	}

	// Join lại và remove extension từ part cuối
	publicID := strings.Join(remainingParts, "/")

	// Remove extension
	if dotIndex := strings.LastIndex(publicID, "."); dotIndex > 0 {
		publicID = publicID[:dotIndex]
	}

	return publicID
}
