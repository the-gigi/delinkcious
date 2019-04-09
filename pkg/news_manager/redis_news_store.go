package news_manager

import (
	"github.com/go-redis/redis"
	"github.com/pelletier/go-toml"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

const redisMaxPageSize = 10

// RedisNewsStore manages a UserEvents data structure
type RedisNewsStore struct {
	redis *redis.Client
}

func (m *RedisNewsStore) GetNews(username string, startIndex int) (events []*om.LinkManagerEvent, nextIndex int, err error) {
	stop := startIndex + redisMaxPageSize - 1
	result, err := m.redis.LRange(username, int64(startIndex), int64(stop)).Result()
	if err != nil {
		return
	}

	for _, t := range result {
		var event om.LinkManagerEvent
		err = toml.Unmarshal([]byte(t), &event)
		if err != nil {
			return
		}

		events = append(events, &event)
	}

	if len(result) == redisMaxPageSize {
		nextIndex = stop + 1
	} else {
		nextIndex = -1
	}

	return
}

func (m *RedisNewsStore) AddEvent(username string, event *om.LinkManagerEvent) (err error) {
	t, err := toml.Marshal(*event)
	if err != nil {
		return
	}

	err = m.redis.RPush(username, t).Err()
	return
}

func NewRedisNewsStore(address string) (store Store, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // use empty password for simplicity. should come from a secret in production
		DB:       0,  // use default DB
	})

	_, err = client.Ping().Result()
	if err != nil {
		return
	}

	store = &RedisNewsStore{redis: client}
	return
}
