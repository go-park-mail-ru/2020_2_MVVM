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

type NewSessionBuilder struct{}

func (sb *NewSessionBuilder) Build(ctx *gin.Context) sessions.Session {
	return sessions.Default(ctx)
}

func GetCurrentUserId(session sessions.Session, userType string) (id uuid.UUID) {
	userIDStr := session.Get(userType)
	if userIDStr == nil {
		return uuid.Nil
	}
	userID, _ := uuid.Parse(userIDStr.(string))
	return userID
}

func GetCandidateOrEmployer(session sessions.Session) (uuid.UUID, string, error) {
	var id uuid.UUID
	if session == nil {
		return uuid.Nil, "", errors.New(AuthRequiredErr)
	}
	if candId := session.Get(CandID); candId != nil {
		id, _ = uuid.Parse(candId.(string))
		return id, Candidate, nil
	} else if emplId := session.Get(EmplID); emplId != nil {
		id, _ = uuid.Parse(emplId.(string))
		return id, Employer, nil
	}
	return uuid.Nil, "", errors.New(AuthRequiredErr)
}
