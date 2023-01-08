package chap_02

import (
	"testing"
	"time"

	"app/utils"

	"app/chap_02/common"
	clientRepo "app/chap_02/repo/client"
	redisConn "app/connector/redis"

	"github.com/google/uuid"
)

func TestLoginCookies(t *testing.T) {
	connector := redisConn.NewRedisConnector()

	client := clientRepo.NewClientRepo(connector)
	token := uuid.New().String()

	t.Run("Test Login cookies", func(t *testing.T) {
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
}
