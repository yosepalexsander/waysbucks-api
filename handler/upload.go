package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
)

type responsePayload struct {
	Filename string `json:"filename"`
}

type response struct {
	commonResponse
	Payload responsePayload `json:"payload"`
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	folder := "waysbucks/avatars"

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		badRequest(w, "maximum upload size is 5 MB")
		return
	}

	file, header, fileErr := r.FormFile("file")
	if fileErr != nil {
		badRequest(w, fileErr.Error())
		return
	}
	defer file.Close()

	if err := helper.ValidateImageFile(header.Header.Get("Content-Type")); err != nil {
		badRequest(w, "upload only for image")
		return
	}

	filename, err := thirdparty.UploadFile(ctx, file, folder, header.Filename)
	if err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(response{commonResponse: commonResponse{
		Message: "resource has successfully created",
	},
		Payload: responsePayload{
			Filename: filename,
		},
	})

	responseOK(w, resp)
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	folder := "waysbucks"

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		badRequest(w, "maximum upload size is 5 MB")
		return
	}

	file, header, fileErr := r.FormFile("file")
	if fileErr != nil {
		badRequest(w, fileErr.Error())
		return
	}
	defer file.Close()

	if err := helper.ValidateImageFile(header.Header.Get("Content-Type")); err != nil {
		badRequest(w, "upload only for image")
		return
	}

	filename, err := thirdparty.UploadFile(ctx, file, folder, header.Filename)
	if err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(response{commonResponse: commonResponse{
		Message: "resource has successfully created",
	},
		Payload: responsePayload{
			Filename: filename,
		},
	})

	responseOK(w, resp)
}
