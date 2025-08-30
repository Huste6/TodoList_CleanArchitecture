package ginupload

import (
	"context"
	"g09/common"
	"g09/component/uploadprovider"
	uploadbiz "g09/module/upload/biz"
	"io"
	"os"
	"path/filepath"
	"strings"

	"mime/multipart"

	"github.com/gin-gonic/gin"
)

// ValidateImageFile validates if the uploaded file is a valid image
func ValidateImageFile(header *multipart.FileHeader) error {
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff"}

	for _, ext := range allowedTypes {
		if fileExt == ext {
			return nil
		}
	}

	return common.NewCustomError(nil, "Unsupported file type. Only images are allowed", "ErrInvalidFileType")
}

// ProcessFileUpload handles the complete file upload process
func ProcessFileUpload(c *gin.Context, file multipart.File, header *multipart.FileHeader) (*common.Image, error) {
	// Validate file type
	if err := ValidateImageFile(header); err != nil {
		return nil, err
	}

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	// Setup upload provider
	provider := uploadprovider.NewCloudinaryProvider(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
		os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),
	)

	// Upload image
	uploadBiz := uploadbiz.NewUploadBiz(provider)
	img, err := uploadBiz.Upload(c.Request.Context(), data, "uploads", header.Filename)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// DeleteOldImage deletes old image from Cloudinary
func DeleteOldImage(ctx context.Context, oldImageURL string) error {
	if oldImageURL == "" {
		return nil
	}

	provider := uploadprovider.NewCloudinaryProvider(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
		os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),
	)

	uploadBiz := uploadbiz.NewUploadBiz(provider)
	return uploadBiz.DeleteImage(ctx, oldImageURL)
}
