package common

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"strings"
)

const (
	ImgDir          = "static"
	MaxImgSize      = 32 << 13 // ~262 Kb
	MaxImgHeight    = 1250     //px
	MaxImgWidth     = 1250     //px
	UploadFileError = -1
	//FileValid       = 0
	PngMime         = "image/png"
	JpegMime        = "image/jpeg"
	base64pngTitle  = 22
	base64jpegTitle = 23
)

func AddOrUpdateUserFile(data io.Reader, imgName string) *Err {
	if data == nil {
		return nil
	}
	fileDir, _ := os.Getwd()
	imgPath := path.Join(fileDir, ImgDir, imgName)

	dst, err := os.Create(imgPath)
	if err != nil {
		return &Err{UploadFileError, "something went wrong (path creation)", nil}
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		return &Err{UploadFileError, "something went wrong (file copy)", nil}
	}
	return nil
}

func imgValidation(imgReader *strings.Reader) *Err {
	if imgReader.Size() > MaxImgSize {
		return &Err{UploadFileError, fmt.Sprintf("file size too big, max allowed: %d kB", MaxImgSize/1000), nil}
	}
	img, _, errDecode := image.DecodeConfig(imgReader)
	if _, errSeek := imgReader.Seek(0, 0); errSeek != nil || errDecode != nil {
		return &Err{UploadFileError, "something went wrong", nil}
	}
	if img.Height > MaxImgHeight || img.Width > MaxImgWidth {
		return &Err{UploadFileError, fmt.Sprintf("the image size exceeds the allowed height %dpx and width %dpx", MaxImgHeight, MaxImgWidth), nil}
	}
	return nil
}

func GetImageFromBase64(data string) (io.Reader, *Err) {
	var imgBase64 string

	if data == "" {
		return nil, nil
	}
	if strings.HasPrefix(data, fmt.Sprintf("data:%s", JpegMime)) && len(data) > base64jpegTitle {
		imgBase64 = data[base64jpegTitle:]
	} else if strings.HasPrefix(data, fmt.Sprintf("data:%s", PngMime)) && len(data) > base64pngTitle {
		imgBase64 = data[base64pngTitle:]
	} else {
		return nil, &Err{code: UploadFileError, message: "unsupported image format, allowed: png/jpeg"}
	}
	imgCode, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return nil, &Err{code: UploadFileError, message: "something wrong with your image, try another one"}
	}
	imgReader := strings.NewReader(string(imgCode))
	if err := imgValidation(imgReader); err != nil {
		return nil, err
	}
	return imgReader, nil
}
