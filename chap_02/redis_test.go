package chap_02

import (
	"testing"
	"time"

	"app/utils"

	"app/chap_02/common"
	clientRepo "app/chap_02/repo/client"
	inventoryCacheRepo "app/chap_02/repo/inventory_cache"
	redisConn "app/connector/redis"

	"github.com/google/uuid"
)

func TestLoginCookies(t *testing.T) {
	connector := redisConn.NewRedisConnector()

	client := clientRepo.NewRepoManager(connector)
	inventory := inventoryCacheRepo.NewRepoManager(connector)

	token := uuid.New().String()

	t.Run("login_cookies", func(t *testing.T) {
		client.UpdateToken(token, "duongnln", "itemX")
		t.Log("Logged-in/update token: ", token)
		t.Log("For user: ", "duongnln")
		t.Log("What duongnln do we get when we look-up that token?")
		r := client.CheckToken(token)
		t.Log("username: ", r)

		utils.AssertStringResult(t, "duongnln", r)

		t.Log("Let's drop the maximum number of cookies to 0 to clean them out")
		t.Log("We will start a thread to do the cleaning, while we stop it later")

		common.LIMIT = 1
		go client.CleanUpSession()
		time.Sleep(1 * time.Second)
		common.QUIT = true
		time.Sleep(2 * time.Second)

		utils.AssertThread(t, common.FLAG)

		s := connector.Conn.HLen("login:").Val()
		t.Log("The current number of sessions still available is:", s)
		utils.AssertnumResult(t, 1, s)
		defer client.Reset()
	})

	t.Run("shopping_cart_cookies", func(t *testing.T) {
		t.Log("Refresh session...")
		client.UpdateToken(token, "duongnln", "itemX")
		t.Log("And add an item to the shopping cart")
		client.AddToCart(token, "itemY", 3)
		client.AddToCart(token, "itemZ", 2)
		r := connector.Conn.HGetAll("cart:" + token).Val()
		t.Log("Our shopping cart currently has:", r)

		utils.AssertTrue(t, len(r) >= 1)

		t.Log("Let's clean out our sessions and carts")
		common.LIMIT = 1
		go client.CleanUpFullSession()
		time.Sleep(1 * time.Second)
		common.QUIT = true
		time.Sleep(2 * time.Second)
		utils.AssertThread(t, common.FLAG)

		r = connector.Conn.HGetAll("cart:" + token).Val()
		t.Log("Our shopping cart now contains:", r)
		defer client.Reset()
	})

	t.Run("Test cache row", func(t *testing.T) {
		t.Log("First, let's schedule caching of itemX every 5 seconds")
		inventory.ScheduleRowCache("itemX", 5)
		t.Log("Our schedule looks like:")
		s := inventory.Connector.Conn.ZRangeWithScores("schedule:", 0, -1).Val()
		t.Log(s[0].Member, s[0].Score)
		utils.AssertTrue(t, len(s) != 0)

		t.Log("We'll start a caching thread that will cache the data...")
		go inventory.CacheRows()
		time.Sleep(1 * time.Second)
		t.Log("Our cached data looks like:")
		r := inventory.Connector.Conn.Get("inv:itemX").Val()
		t.Log(r)
		utils.AssertTrue(t, len(r) != 0)

		t.Log("We'll check again in 5 seconds...")
		time.Sleep(5 * time.Second)
		t.Log("Notice that the data has changed...")
		r2 := inventory.Connector.Conn.Get("inv:itemX").Val()
		t.Log(r2)
		utils.AssertTrue(t, len(r) != 0)
		utils.AssertTrue(t, r != r2)

		t.Log("Let's force un-caching")
		inventory.ScheduleRowCache("itemX", -1)
		time.Sleep(1 * time.Second)
		r = inventory.Connector.Conn.Get("inv:itemX").Val()
		t.Log("The cache was cleared?", r == "")
		utils.AssertFalse(t, r != "")

		common.QUIT = true
		time.Sleep(2 * time.Second)
		utils.AssertThread(t, common.FLAG)
		defer inventory.Connector.Conn.FlushDB()
	})
}
