package thirdparty

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var (
	ErrServiceUnavailable error = errors.New("object storage service unavailable")
)

func UploadFile(ctx context.Context, file multipart.File, filename string) error {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		fmt.Printf("Failed to intialize Cloudinary, %v", err)
		return err
	}

	// Upload file image to Cloudinary
	_, uploadErr := cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{PublicID: filename, UniqueFilename: true, UseFilename: true},
	)

	if uploadErr != nil {
		fmt.Println(uploadErr)
		return err
	}
	return nil
}

func GetImageUrl(ctx context.Context, publicID string) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		fmt.Printf("Failed to intialize Cloudinary, %v", err)
		return "", err
	}

	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: publicID})

	if err != nil {
		log.Printf("err %v", err)
		return "", err
	}
	return asset.SecureURL, nil
}

func RemoveFile(ctx context.Context, filename string) error {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		fmt.Printf("Failed to intialize Cloudinary, %v", err)
		return err
	}

	if _, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: filename}); err != nil {
		return err
	}
	return nil
}
