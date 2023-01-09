package inventory_cache

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"app/chap_02/common"
	"app/chap_02/model"

	redisConn "app/connector/redis"

	"github.com/go-redis/redis"
)

func getOneByID(id string) model.Inventory {
	return model.NewInventory(id, "data to cache...", time.Now().Unix())
}

type repoManager struct {
	Connector *redisConn.RedisConnector
}

func NewRepoManager(connector *redisConn.RedisConnector) *repoManager {
	return &repoManager{
		Connector: connector,
	}
}

func (r *repoManager) ScheduleRowCache(rowId string, delay int64) {
	r.Connector.Conn.ZAdd("delay:", redis.Z{Member: rowId, Score: float64(delay)})
	r.Connector.Conn.ZAdd("schedule:", redis.Z{Member: rowId, Score: float64(time.Now().Unix())})
}

func (r *repoManager) CacheRows() {
	for !common.QUIT {
		next := r.Connector.Conn.ZRangeWithScores("schedule:", 0, 0).Val()
		now := time.Now().Unix()
		if len(next) == 0 || next[0].Score > float64(now) {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		rowId := next[0].Member.(string)
		delay := r.Connector.Conn.ZScore("delay:", rowId).Val()
		if delay <= 0 {
			r.Connector.Conn.ZRem("delay:", rowId)
			r.Connector.Conn.ZRem("schedule:", rowId)
			r.Connector.Conn.Del("inv:" + rowId)
			continue
		}

		row := getOneByID(rowId)
		r.Connector.Conn.ZAdd("schedule:", redis.Z{Member: rowId, Score: float64(now) + delay})
		jsonRow, err := json.Marshal(row)
		if err != nil {
			log.Fatalf("marshal json failed, data is: %v, err is: %v\n", row, err)
		}
		r.Connector.Conn.Set("inv:"+rowId, jsonRow, 0)
	}
	defer atomic.AddInt32(&common.FLAG, -1)
}
