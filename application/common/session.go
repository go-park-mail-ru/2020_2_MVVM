package common

import (
	"errors"
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

func GetCurrentUserId(session sessions.Session, userType string) (id uuid.UUID, err error) {
	userIDStr := session.Get(userType)
	if userIDStr == nil {
		return uuid.Nil, nil
	}
	userID, _ := uuid.Parse(userIDStr.(string))
	return userID, nil
}

func GetUser(session sessions.Session) (uuid.UUID, string, error) {   // <3
	var id uuid.UUID
	if userId := session.Get(UserID); userId != nil {
		id, _ = uuid.Parse(userId.(string))
		return id, User, nil
	} else if emplId := session.Get(EmplID); emplId != nil {
		id, _ = uuid.Parse(emplId.(string))
		return id, Employer, nil
	} else if candId := session.Get(CandID); candId != nil {
		id, _ = uuid.Parse(candId.(string))
		return id, Candidate, nil
	}
	return uuid.Nil, "", errors.New(AuthRequiredErr)
}
