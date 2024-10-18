package abtest

import (
	"context"
	"encoding/json"
	"github.com/spaolacci/murmur3"
	"recommend-server/internal/common/middleware"
)

type ExpConfig struct {
	Layer     string      `json:"layer"`      //层
	Seed      uint32      `json:"seed"`       //随机种子
	ExpGroups []*ExpGroup `json:"exp_groups"` //实验组   a:a的参数  b:b的参数
}

type ExpGroup struct {
	Name      string        `json:"name"`
	Start     uint64        `json:"start"`
	End       uint64        `json:"end"`
	Parameter *ExpParameter `json:"parameter"`
}

type ExpParameter struct {
	RecallConfig RecallConfig `json:"recall_config"`
	MergeConfig  MergeConfig  `json:"merge_config"`
	RankConfig   RankConfig   `json:"rank_config"`
	ReRankConfig ReRankConfig `json:"rerank_config"`
	FilterConfig FilterConfig `json:"filter_config"`
}

type RecallConfig struct {
	IndexId        string                 `json:"index_id"`
	RecallStrategy string                 `json:"recall_strategy"`
	Param          map[string]interface{} `json:"param"`
}

type MergeConfig struct {
	FuncName string                 `json:"func_name"`
	Param    map[string]interface{} `json:"param"`
}

type RankConfig struct {
	FuncName string                 `json:"func_name"`
	Param    map[string]interface{} `json:"param"`
}

type ReRankConfig struct {
	FuncName string                 `json:"func_name"`
	Param    map[string]interface{} `json:"param"`
}

type FilterConfig struct {
	FuncName string                 `json:"func_name"`
	Param    map[string]interface{} `json:"param"`
}

func GenerateHashNum(group int, seed int, hashKey string) uint64 {
	h64 := murmur3.New64WithSeed(uint32(seed))
	h64.Write([]byte(hashKey))
	hashInt := h64.Sum64()
	return hashInt % uint64(group)
}

func GetExpConfig(redisKey string) *ExpConfig {
	expConfig := &ExpConfig{}
	if middleware.M.RedisCli != nil {
		values, _ := middleware.M.RedisCli.Get(context.Background(), redisKey).Result()
		err := json.Unmarshal([]byte(values), expConfig)
		if err != nil {
			return nil
		}
		return expConfig
	} else if middleware.M.RedisClusterCli != nil {
		values, _ := middleware.M.RedisClusterCli.Get(context.Background(), redisKey).Result()
		err := json.Unmarshal([]byte(values), expConfig)
		if err != nil {
			return nil
		}
		return expConfig
	} else {
		return nil
	}
}

type WhiteList struct {
	Wl map[string]string `json:"wl"`
}

func GetWhiteList(redisKey string) *WhiteList {
	whiteList := &WhiteList{}
	if middleware.M.RedisCli != nil {
		values, _ := middleware.M.RedisCli.Get(context.Background(), redisKey).Result()
		err := json.Unmarshal([]byte(values), whiteList)
		if err != nil {
			return nil
		}
		return whiteList
	} else if middleware.M.RedisClusterCli != nil {
		values, _ := middleware.M.RedisClusterCli.Get(context.Background(), redisKey).Result()
		err := json.Unmarshal([]byte(values), whiteList)
		if err != nil {
			return nil
		}
		return whiteList
	} else {
		return nil
	}
}
