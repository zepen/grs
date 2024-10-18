package abtest

import (
	"fmt"
	"math/rand"
	"recommend-server/internal/common/middleware"
	"testing"
	"time"
)

func generateID() string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	randomNumber := rand.Intn(1000000000) // 生成一个小于 1000000000 的随机数
	// 将时间戳和随机数拼接成 19 位数的 ID
	id := fmt.Sprintf("%013d%09d", timestamp, randomNumber)
	return id
}

func TestHash(t *testing.T) {
	recKey := "exp:abtest:recommend"
	j := 0
	group := 100
	redisConf := &middleware.RedisConf{}
	redisConf.Host = "127.0.0.1:6379"
	middleware.M = &middleware.MiddleWare{}
	middleware.M.RedisCli = middleware.NewRedisClient(redisConf)
	expConfig := GetExpConfig(recKey)
	if expConfig != nil {
		expGroups := group / len(expConfig.ExpGroups)
		fmt.Printf("exp len = %d\n", expGroups)
		fmt.Printf("--------------------\n")
		for j < 10 {
			mp := make(map[string]int)
			seed := int(expConfig.Seed)
			i := 0
			all := 1000
			for i < all {
				hashKey := generateID()
				bb := GenerateHashNum(group, seed, hashKey)
				for k, expKey := range expConfig.ExpGroups {
					// fmt.Printf("%d, start:%d, end:%d;\n", bb, expGroups*k, expGroups*(k+1))
					if bb >= uint64(expGroups*k) && bb < uint64(expGroups*(k+1)) {
						mp[expKey.Name]++
					}
				}
				i++
			}
			for _, expKey := range expConfig.ExpGroups {
				fmt.Printf("%s:%d:%f, start:%d,end:%d,recall_strategy:%s\n",
					expKey.Name,
					mp[expKey.Name],
					float32(mp[expKey.Name])/float32(all),
					expKey.Start,
					expKey.End,
					expKey.Parameter.RecallConfig.RecallStrategy,
				)
			}
			fmt.Printf("--------------------\n")
			j++
		}
	} else {
		fmt.Printf("expConfig is nil!\n")
	}
}

func TestWhiteList(t *testing.T) {
	recKey := "exp:abtest:recommend:whitelist"
	redisConf := &middleware.RedisConf{}
	redisConf.Host = "127.0.0.1:6379"
	middleware.M = &middleware.MiddleWare{}
	middleware.M.RedisCli = middleware.NewRedisClient(redisConf)
	fmt.Println(GetWhiteList(recKey).Wl)
}
