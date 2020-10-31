package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	ImgDir          = "static"
	MaxImgSize      = 32 << 13 // ~262 Kb
	MaxImgHeight    = 1250      //px
	MaxImgWidth     = 1250      //px
	UploadFileError = -1
	FileValid       = 0
	EmptyField		= 1
)

func fileValidation(header *multipart.FileHeader, file multipart.File, allowedFormats []string, allowedSize int64) Err {
	if header == nil {
		return NewErr(UploadFileError, "something went wrong", nil)
	}
	fileType := header.Header.Get("Content-Type")
	extWasFind := false
	for i := 0; i < len(allowedFormats); i++ {
		if fileType == allowedFormats[i] {
			extWasFind = true
		}
	}
	if extWasFind == false {
		return NewErr(UploadFileError, "not supported file extension", allowedFormats)
	}
	if header.Size > allowedSize {
		return NewErr(UploadFileError, fmt.Sprintf("file size too big, max allowed: %d kB", allowedSize/1000), nil)
	}
	//TODO: if in project will be some file upload like resume or certificates move validation by format in different functions
	// Image validation(avatar)
	img, _, errDecode := image.DecodeConfig(file)
	if _, errSeek := file.Seek(0, 0); errSeek != nil || errDecode != nil {
		return NewErr(UploadFileError, "something went wrong", nil)
	}
	if img.Height > MaxImgHeight || img.Width > MaxImgWidth {
		return NewErr(UploadFileError, fmt.Sprintf("the image size exceeds the allowed height %dpx and width %dpx", MaxImgHeight, MaxImgWidth), nil)
	}
	return NewErr(FileValid, "", nil)
}

func AddOrUpdateUserImage(data io.Reader, imgPath string) error {
	path := filepath.Join(ImgDir, imgPath)

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		return err
	}
	return nil
}

func GetImage(ctx *gin.Context, imgName string) (*multipart.File, *Err) {
	file, header, err := ctx.Request.FormFile(imgName)
	if err == nil {
		if err := fileValidation(header, file, []string{"image/jpeg", "image/png"}, MaxImgSize); err.Code() != FileValid {
			return nil, &Err{code:http.StatusOK, message: err.String()}
		}
	} else if err.Error() != "http: no such file" {
		return nil, &Err{code: http.StatusBadRequest, message: err.Error()}
	} else {
		return nil, &Err{code: EmptyField}
	}
	return &file, nil
}
