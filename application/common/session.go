package common

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionBuilder interface {
	Build(ctx *gin.Context) sessions.Session
}

type NewSessionBuilder struct {}

func (sb *NewSessionBuilder) Build (ctx *gin.Context) sessions.Session {
	return sessions.Default(ctx)
}

func GetCurrentUserId(ctx *gin.Context, user_type string) (id uuid.UUID, err error) {
	session := sessions.Default(ctx)
	userIDStr := session.Get(user_type)
	if userIDStr == nil {
		return uuid.Nil, nil
	}
	userID, err := uuid.Parse(userIDStr.(string))

	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}


