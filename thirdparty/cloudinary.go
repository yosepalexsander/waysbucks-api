package thirdparty

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"golang.org/x/sync/errgroup"
)

var (
	ErrServiceUnavailable error = errors.New("object storage service unavailable")
)

func UploadFile(ctx context.Context, file multipart.File, filename string) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		fmt.Printf("Failed to intialize Cloudinary\nerror: %v", err)
		return "", ErrServiceUnavailable
	}

	name := strings.Split(filename, ".")[0] + "-" + helper.RandString(15)

	_, err = cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{PublicID: name, Format: filepath.Ext(filename)[1:]},
	)

	if err != nil {
		fmt.Printf("Failed to upload file\nerror: %v", err)
		return "", err
	}

	return name, nil
}

func GetImageUrl(ctx context.Context, publicID string) (string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		fmt.Printf("Failed to intialize Cloudinary, %v", err)
		return "", ErrServiceUnavailable
	}

	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: publicID})
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	return asset.SecureURL, nil
}

func RemoveFile(ctx context.Context, filename string) error {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		fmt.Printf("Failed to intialize Cloudinary, %v", err)
		return ErrServiceUnavailable
	}

	if _, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: filename}); err != nil {
		return err
	}

	return nil
}

func UpdateImage(file multipart.File, oldName string, newName string) (string, error) {
	filename := ""
	g, ctx := errgroup.WithContext(context.TODO())

	g.Go(func() error {
		var err error
		if filename, err = UploadFile(ctx, file, newName); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := RemoveFile(ctx, oldName); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	return filename, nil
}
