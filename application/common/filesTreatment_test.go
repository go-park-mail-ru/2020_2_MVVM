package common

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetImageFromBase64(t *testing.T) {
	testImg := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString([]byte("testImg")))
	readerNil, err := GetImageFromBase64("")
	assert.Nil(t, readerNil)
	assert.Nil(t, err)
	readerNil, err = GetImageFromBase64("testImg")
	assert.Nil(t, readerNil)
	assert.Error(t, err)
	validImgPrefix1 := fmt.Sprintf("data:%s,someImage", JpegMime)
	readerNil, err = GetImageFromBase64(validImgPrefix1)
	assert.Nil(t, readerNil)
	assert.Error(t, err)
	validImgPrefix2 := fmt.Sprintf("data:%s,someImage", PngMime)
	readerNil, err = GetImageFromBase64(validImgPrefix2)
	assert.Nil(t, readerNil)
	assert.Error(t, err)
	readerNil, err = GetImageFromBase64(testImg)
	assert.Nil(t, readerNil)
	assert.Error(t, err)
}

func TestAddOrUpdateUserFile(t *testing.T) {
	err := AddOrUpdateUserFile(nil, "someImgName")
	assert.Nil(t, err)
	err = AddOrUpdateUserFile(strings.NewReader("temp"), ".WrongPath")
	assert.Error(t, err)
}