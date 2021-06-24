package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type AccessDetails struct {
	TokenUuid string
	UserId    string
	UserName  string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type IAuth interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(string) (string, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type RedisAuthService struct {
	client *redis.Client
}

var _ IAuth = &RedisAuthService{}

func NewAuthService(client *redis.Client) *RedisAuthService {
	return &RedisAuthService{client: client}
}

func (tk *RedisAuthService) CreateAuth(userId string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(td.TokenUuid, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := tk.client.Set(td.RefreshUuid, userId, rt.Sub(now)).Result()
	if err != nil {
		return err
	}

	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}

	return nil
}

func (tk *RedisAuthService) FetchAuth(tokenUuid string) (string, error) {
	userid, err := tk.client.Get(tokenUuid).Result()
	if err != nil {
		return "", err
	}

	return userid, nil
}

func (tk *RedisAuthService) DeleteTokens(authD *AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)
	deletedAt, err := tk.client.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}

	deletedRt, err := tk.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}

	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}

	return nil
}

func (tk *RedisAuthService) DeleteRefresh(refreshUuid string) error {
	deleted, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}

	return nil
}
