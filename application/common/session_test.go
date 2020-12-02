package common

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

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