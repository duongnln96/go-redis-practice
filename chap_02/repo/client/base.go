package client

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"app/chap_02/common"
	redisConn "app/connector/redis"
	"app/utils"

	"github.com/go-redis/redis"
)

type clientRepo struct {
	Connector *redisConn.RedisConnector
}

func NewClientRepo(connector *redisConn.RedisConnector) *clientRepo {
	return &clientRepo{
		Connector: connector,
	}
}

func (r *clientRepo) Reset() {
	r.Connector.Conn.FlushDB()

	common.QUIT = false
	common.LIMIT = 10000000
	common.FLAG = 1
}

func (r *clientRepo) CheckToken(token string) string {
	return r.Connector.Conn.HGet("login:", token).Val()
}

// UpdateToken
//
// login - hash - key token value user
//
// zset contains score and member (score is used to sort)
// recent - set sort by timestamp
//
// view - set sort by timestamp
func (r *clientRepo) UpdateToken(token, user, item string) {
	timestamp := time.Now().Unix()
	r.Connector.Conn.HSet("login:", token, user)
	r.Connector.Conn.ZAdd("recent:", redis.Z{
		Score:  float64(timestamp),
		Member: token,
	})

	if item != "" {
		r.Connector.Conn.ZAdd(fmt.Sprintf("viewed:%s", token), redis.Z{
			Score:  float64(timestamp),
			Member: item,
		})
		r.Connector.Conn.ZRemRangeByRank(fmt.Sprintf("viewed:%s", token), 0, -26)
	}
}

// Over time, memory use will grow, and we’ll want to clean out old data
func (r *clientRepo) CleanUpSession() {
	for !common.QUIT {
		size := r.Connector.Conn.ZCard("recent:").Val()
		if size <= common.LIMIT {
			time.Sleep(1 * time.Second)
			continue
		}

		endIndex := utils.Min(size-common.LIMIT, 100)
		log.Printf("===endIndex===: %d", endIndex)
		tokens := r.Connector.Conn.ZRange("recent:", 0, endIndex-1).Val()
		log.Printf("===tokens===: %+v", tokens)

		var sessionKeys []string = make([]string, 0, len(tokens))
		for _, token := range tokens {
			sessionKeys = append(sessionKeys, fmt.Sprintf("viewed:%s", token))
		}

		r.Connector.Conn.Del(sessionKeys...)
		r.Connector.Conn.HDel("login:", tokens...)
		r.Connector.Conn.ZRem("recent:", tokens)
	}

	defer atomic.AddInt32(&common.FLAG, -1)
}

// cart is hash table with name cart:token contains key is item and value is count
func (r *clientRepo) AddToCart(session, item string, count int) {
	if count <= 0 {
		r.Connector.Conn.HDel(fmt.Sprintf("cart:%s", session), item)
	} else {
		r.Connector.Conn.HSet(fmt.Sprintf("cart:%s", session), item, count)
	}
}

// Over time, memory use will grow, and we’ll want to clean out old data
func (r *clientRepo) CleanUpFullSession() {
	for !common.QUIT {
		size := r.Connector.Conn.ZCard("recent:").Val()
		if size <= common.LIMIT {
			time.Sleep(1 * time.Second)
			continue
		}

		endIndex := utils.Min(size-common.LIMIT, 100)
		// log.Printf("===endIndex===: %d", endIndex)
		tokens := r.Connector.Conn.ZRange("recent:", 0, endIndex-1).Val()
		// log.Printf("===tokens===: %+v", tokens)

		var sessionKeys []string = make([]string, 0, len(tokens))
		for _, token := range tokens {
			sessionKeys = append(sessionKeys, fmt.Sprintf("viewed:%s", token))
			sessionKeys = append(sessionKeys, fmt.Sprintf("cart:%s", token))
		}

		r.Connector.Conn.Del(sessionKeys...)
		r.Connector.Conn.HDel("login:", sessionKeys...)
		r.Connector.Conn.ZRem("recent:", sessionKeys)
	}

	defer atomic.AddInt32(&common.FLAG, -1)
}
