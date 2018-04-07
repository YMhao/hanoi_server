package dao

import (
	"strings"

	"github.com/mediocregopher/radix.v2/redis"
	"github.com/satori/go.uuid"
)

type _SessionDao struct {
	client *redis.Client
}

func newSessionDao(client *redis.Client) *_SessionDao {
	return &_SessionDao{
		client: client,
	}
}

func (d _SessionDao) NewSession(userUUID string) (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionId := "S" + strings.Replace(uuid.String(), "-", "", -1)

	err = d.client.Cmd("SET", sessionId, userUUID, "EX", 3600).Err
	if err != nil {
		return "", nil
	}
	return sessionId, nil
}

func (d _SessionDao) GetUserUUID(sessionId string) (string, error) {
	return d.client.Cmd("GET", sessionId).Str()
}
