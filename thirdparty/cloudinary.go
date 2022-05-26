package thirdparty

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/yosepalexsander/waysbucks-api/helper"
)

var (
	ErrServiceUnavailable error = errors.New("object storage service unavailable")
)

func UploadFile(ctx context.Context, file multipart.File, filename string, folder string) (string, string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Printf("Failed to intialize Cloudinary\nerror: %v", err)
		return "", "", ErrServiceUnavailable
	}

	publicID := strings.Split(filename, ".")[0] + "-" + helper.RandString(15)
	format := filepath.Ext(filename)

	uploadResult, err := cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{PublicID: publicID, Folder: folder, Format: format[1:]},
	)

	if err != nil {
		log.Printf("Failed to upload file\nerror: %v", err)
		return "", "", err
	}

	return uploadResult.PublicID, uploadResult.SecureURL, nil
}

func GetImageUrl(ctx context.Context, publicID string) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Printf("Failed to intialize Cloudinary, %v", err)
		return "", ErrServiceUnavailable
	}

	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: publicID})
	if err != nil {
		log.Printf("Failed to retrieve asset details\nerror: %v", err)
		return "", err
	}

	return asset.SecureURL, nil
}

func RemoveFile(ctx context.Context, publicID string) error {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Printf("Failed to intialize Cloudinary, %v", err)
		return ErrServiceUnavailable
	}

	if _, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID}); err != nil {
		return err
	}

	return nil
}
