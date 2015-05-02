package queue

import (
	"time"

	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

const batchSize = 250
const timeout = 5

func InitQueue(redisServer string) {
	pool = newPool(redisServer)
}

func FetchBatch() <-chan []string {
	return listFetch("jobs")
}

func FetchResults() <-chan []string {
	return listFetch("results")
}

func FetchNotifications() <-chan []string {
	return listFetch("notifications")
}

func WriteResult(result string) {
	c := pool.Get()
	defer c.Close()

	_, err := redis.String(c.Do("LPUSH", "results", result))
	if err != nil {
		// never mind
	}

	return
}

func AcquireSchedulerLock(interval, timeout int) bool {
	c := pool.Get()
	defer c.Close()

	exists, err := redis.Bool(c.Do("EXISTS", "scheduler:last_run"))
	if err != nil || exists {
		return false
	}

	_, err = redis.String(c.Do("SET", "scheduler:lock", "LOCK", "NX", "EX", timeout))
	if err != nil {
		return false
	}

	c.Do("SET", "scheduler:last_run", time.Now().Format(time.RFC3339), "EX", interval)

	return true
}

func ReleaseSchedulerLock() {
	c := pool.Get()
	defer c.Close()
	c.Do("DEL", "scheduler:lock")
}

func listWrite(list, content string) {
	c := pool.Get()
	defer c.Close()

	c.Do("LPUSH", list, content)
}

func AddJob(job string) {
	listWrite("jobs", job)
}

func AddNotification(notification string) {
	listWrite("notifications", notification)
}

func listFetch(listName string) <-chan []string {
	out := make(chan []string)

	go func() {
		c := pool.Get()
		defer c.Close()

		results := []string{}

		for i := 0; i < batchSize; i++ {
			reply, err := redis.Values(c.Do("BRPOP", listName, timeout))
			if err != nil {
				break
			}
			item := string(reply[1].([]byte))
			results = append(results, item)
		}

		out <- results
	}()

	return out
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
