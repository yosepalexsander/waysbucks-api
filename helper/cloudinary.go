package helper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"regexp"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var (
	ErrorInvalidFileExtension error = errors.New("invalid file extension")
)
func UploadFile(ctx context.Context, file multipart.File, filename string) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		fmt.Printf("Failed to intialize Cloudinary, %v", err)
		return "", err
	}

	regex, _ := regexp.Compile(`\.(jpg|JPEG|png|PNG|svg|SVG)$`)

	// Check if file extension match the regex or not
	if isMatch := regex.MatchString(filename); !isMatch {
		return "", ErrorInvalidFileExtension
	}

	// Split filename and file extension
	filename = strings.Split(filename, ".")[0] + RandString(20)

	// Upload file image to Cloudinary
	uploadResult, uploadErr := cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{ PublicID: filename, UniqueFilename: true, UseFilename: true },
	)

	if uploadErr != nil {
		fmt.Println(uploadErr)
		return "", err
	}
	return uploadResult.PublicID, nil
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