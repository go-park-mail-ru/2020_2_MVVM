package common

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func FileValidation(header *multipart.FileHeader, allowedFormats []string, allowedSize int64) error {
	fileExt := filepath.Ext(header.Filename)
	extWasFind := false
	for i := 0; i < len(allowedFormats); i++ {
		if fileExt == allowedFormats[i] {
			extWasFind = true
		}
	}
	if extWasFind == false || (fileExt != ".png" && fileExt != ".jpeg") {
		return NewErr(uploadImgError, "not supported file extension", allowedFormats)
	}
	if header.Size > allowedSize {
		return errors.New(fmt.Sprintf("file size too big, max allowed: %d", allowedSize))
	}
	return nil
}

func addOrUpdateUserImage(imgPath string, data io.Reader) error {
	path := filepath.Join(imgDir, imgPath)

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
