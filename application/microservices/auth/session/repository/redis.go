package repository

import (
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/session"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisSessions struct {
	conn            *redis.Client
	ctx             context.Context
	logger          common.Logger
	sessionDuration time.Duration
}

func NewRedisSessionRepository(conn *redis.Client, log common.Logger) session.Repository {
	return &redisSessions{conn: conn, ctx: context.Background(), sessionDuration: time.Hour * 24 * 32, logger: log}
}

func (sd *redisSessions) Add(sessionID string, session common.BasicSession) error {
	p, err := json.Marshal(session)
	if err != nil {
		sd.logger.Error.Error(err.Error())
		return err
	}
	err = sd.conn.Set(sd.ctx, sessionID, p, sd.sessionDuration).Err()
	if err != nil {
		sd.logger.Error.Error(err.Error())
	}
	return nil
}

func (sd *redisSessions) GetSession(sessionID string) (session *common.BasicSession, err error) {
	p, err := sd.conn.Get(sd.ctx, sessionID).Result()
	if err != nil {
		sd.logger.Error.Error(err.Error())
		return nil, err
	}

	session = new(common.BasicSession)
	err = json.Unmarshal([]byte(p), session)

	if err != nil {
		sd.logger.Error.Error(err.Error())
		return nil, err
	}
	return session, nil
}

func (sd *redisSessions) Delete(sessionID string) error {
	err := sd.conn.Del(sd.ctx, sessionID).Err()
	if err != nil {
		sd.logger.Error.Error(err.Error())
	}
	return err
}
