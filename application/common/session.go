package common

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Session interface {
	GetSessionID() string
	GetUserID() uuid.UUID
	GetEmplID() uuid.UUID
	GetCandID() uuid.UUID
}

type SessionBuilder interface {
	Build(ctx *gin.Context) Session
}

type NewSessionBuilder struct{}

func (sb *NewSessionBuilder) Build(ctx *gin.Context) Session {
	session, exist := ctx.Get("session")
	if !exist {
		return nil
	}
	return session.(Session)
}

type BasicSession struct {
	SessionID string
	UserID    uuid.UUID
	EmplID    uuid.UUID
	CandID    uuid.UUID
}

func (s *BasicSession) GetSessionID() string {
	return s.SessionID
}

func (s *BasicSession) GetUserID() uuid.UUID {
	return s.UserID
}

func (s *BasicSession) GetEmplID() uuid.UUID {
	return s.EmplID
}

func (s *BasicSession) GetCandID() uuid.UUID {
	return s.CandID
}
