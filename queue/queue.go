package queue

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var RedisAddressWithPort string

const batchSize = 10
const timeout = 5

func FetchBatch() (result []string) {
	return listFetch("jobs")
}

func FetchResults() (results []string) {
	return listFetch("results")
}

func WriteResult(result string) {
	fmt.Println("Writing result " + result)
	c, err := redis.Dial("tcp", RedisAddressWithPort)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	_, err = redis.String(c.Do("LPUSH", "results", result))
	if err != nil {
		// never mind
	}

	return
}

func AddJob(job string) {
	c, err := redis.Dial("tcp", RedisAddressWithPort)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.Do("LPUSH", "jobs", job)
	fmt.Println("Job added:", job)
}

func listFetch(listName string) (results []string) {
	c, err := redis.Dial("tcp", RedisAddressWithPort)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for i := 0; i < batchSize; i++ {
		reply, err := redis.Values(c.Do("BRPOP", listName, timeout))
		if err != nil {
			break
		}
		item := string(reply[1].([]byte))
		results = append(results, item)
	}

	return results
}
