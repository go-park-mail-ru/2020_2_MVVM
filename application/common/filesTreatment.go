package common

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func FileValidation(header *multipart.FileHeader, allowedFormats []string, allowedSize int64) Err {
	if header == nil {
		return NewErr(UploadFileError, "something went wrong", nil)
	}
	fileExt := filepath.Ext(header.Filename)
	extWasFind := false
	for i := 0; i < len(allowedFormats); i++ {
		if fileExt == allowedFormats[i] {
			extWasFind = true
		}
	}
	if extWasFind == false || (fileExt != ".png" && fileExt != ".jpeg") {
		return NewErr(UploadFileError, "not supported file extension", allowedFormats)
	}
	if header.Size > allowedSize {
		return NewErr(UploadFileError, fmt.Sprintf("file size too big, max allowed: %d kB", allowedSize / 1000), nil)
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
