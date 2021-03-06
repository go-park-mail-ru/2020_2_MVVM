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
	MaxImgSize      = 32 << 16 // 2 Mb
	MaxImgHeight    = 2500     //px
	MaxImgWidth     = 2500     //px
	UploadFileError = -1
	//FileValid       = 0
	PngMime         = "image/png"
	JpegMime        = "image/jpeg"
	base64pngTitle  = 22
	base64jpegTitle = 23
	someWentWrong   = "Что-то пошло не так, попробуйте позже."
)

func AddOrUpdateUserFile(data io.Reader, imgName string) *Err {
	if data == nil {
		return nil
	}
	//fileDir, _ := os.Getwd()
	imgPath := path.Join(PathToSaveStatic, imgName)

	dst, err := os.Create(imgPath)
	if err != nil {
		return &Err{UploadFileError, someWentWrong, nil}
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		return &Err{UploadFileError, someWentWrong, nil}
	}
	return nil
}

func imgValidation(imgReader *strings.Reader) *Err {
	if imgReader.Size() > MaxImgSize {
		return &Err{UploadFileError, fmt.Sprintf("Превышен максимальный размер изображения. Максимальный размер: %d mB.", MaxImgSize/(1024*1024)), nil}
	}
	img, _, errDecode := image.DecodeConfig(imgReader)
	if _, errSeek := imgReader.Seek(0, 0); errSeek != nil || errDecode != nil {
		return &Err{UploadFileError, someWentWrong, nil}
	}
	if img.Height > MaxImgHeight || img.Width > MaxImgWidth {
		return &Err{UploadFileError, fmt.Sprintf("Размеры изображения превышают допутимую высоту %dpx и ширину %dpx.", MaxImgHeight, MaxImgWidth), nil}
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
		return nil, &Err{code: UploadFileError, message: "Неподдерживаемый формат изображения, разрешены любые форматы png/jpeg."}
	}
	imgCode, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return nil, &Err{code: UploadFileError, message: "Что-то не так с вашим изображением, попробуйте выбрать другое."}
	}
	imgReader := strings.NewReader(string(imgCode))
	if err := imgValidation(imgReader); err != nil {
		return nil, err
	}
	return imgReader, nil
}
