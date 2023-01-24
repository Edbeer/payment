package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Edbeer/auth-grpc/types"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)

type RedisStorage struct {
	redis  *redis.Client
}

func NewRedisStorage(redis  *redis.Client) *RedisStorage {
	return &RedisStorage{
		redis: redis,
	}
}

// Add refresh token in redis
func (s *RedisStorage) CreateSession(ctx context.Context, session *types.Session, expire int) (string, error) {
	session.RefreshToken = newRefreshToken()

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return "", err
	}
	if err := s.redis.Set(ctx, session.RefreshToken, sessionBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", err
	}

	return session.RefreshToken, nil
}

// Get user id from session
func (s *RedisStorage) GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	sessionBytes, err := s.redis.Get(ctx, refreshToken).Bytes()
	if err != nil {
		return uuid.Nil , err
	}
	session := &types.Session{}
	if err = json.Unmarshal(sessionBytes, session); err != nil {
		return uuid.Nil, err
	}

	return session.UserID, nil
}

// Delete session cookie
func (s *RedisStorage) DeleteSession(ctx context.Context, refreshToken string) error {
	if err := s.redis.Del(ctx, refreshToken).Err(); err != nil {
		return err
	}
	return nil
}

func newRefreshToken() string {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", b)
}