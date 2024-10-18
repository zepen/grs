package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-redis/redis/v8"
	"log"
	"strings"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redisConf := &RedisConf{}
	redisConf.Host = "127.0.0.1:6379"
	redisCli := NewRedisClient(redisConf)
	if redisCli != nil {
		log.Printf("create redis connection completed!")
	}
	values, _ := redisCli.Get(context.Background(), "user_intent_list:1635093838685372418").Result()
	for i, v := range strings.Split(values, ",") {
		fmt.Println(i, v, strings.Replace(v, "\"", "", -1))
	}
	fmt.Println(values, strings.Replace(values, "\"", "", -1))
}

func TestRedis2(t *testing.T) {
	redisConf := &RedisConf{}
	redisConf.Host = "192.168.31.214:16379"
	redisConf.Password = "cher8"
	// redisConf.Host = "52.15.231.163:6379"
	redisCli := NewRedisClient(redisConf)
	if redisCli != nil {
		log.Printf("create redis connection completed!")
	}
	values, err := redisCli.Get(context.Background(), "AD:CONFIGURATION:823371426739392512").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(values)
	values = strings.Replace(values[1:len(values)-1], "\\", "", -1)
	fmt.Println(values)
	var ad map[string]interface{}
	err2 := json.Unmarshal([]byte(values), &ad)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(time.Now().UnixMilli())
	fmt.Println(ad["startTime"].(float64))
	fmt.Println(ad["endTime"].(float64))
	fmt.Println(ad["limitAmount"].(float64))
}

func TestRedisCluster(t *testing.T) {
	host := "redis-cluster.yx5mgd.use2.cache.amazonaws.com:36379"
	var addrs []string
	for i := 1; i <= 2; i++ {
		for j := 1; j <= 2; j++ {
			addrs = append(addrs, fmt.Sprintf("redis-cluster-000%d-00%d.%s", i, j, host))
		}
	}
	fmt.Println(addrs)
	//clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: addrs,
	//})
	//_, err := clusterClient.Ping(context.TODO()).Result()
	//if err != nil {
	//	log.Printf("[Error] Redis connection is fail: %s", err)
	//}
}
