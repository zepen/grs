package dump

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"gitlab.com/cher8/lion/common/ilog"
	"gitlab.com/cher8/lion/common/middleware"
	"gitlab.com/cher8/lion/common/robot"
	"os"
	"recommend-server/internal/model"
	"strconv"
	"strings"
	"time"
)

func NewRecallTypeDump(reqContext *model.NoteReqContext, note *model.Note) string {
	recallTypeDump := make(map[string]interface{})
	recallTypeDump["rid"] = reqContext.ReqId
	recallTypeDump["user_id"] = reqContext.UserId
	recallTypeDump["note_id"] = note.NoteId
	recallTypeDump["recall_name"] = note.RecallScore.Name
	if reqContext.AbTestConf.Name != "" {
		recallTypeDump["bucket"] = reqContext.AbTestConf.Name
	}
	bb, _ := json.Marshal(recallTypeDump)
	return string(bb)
}

func SendInfo2Kafka(info string) {
	topic := "recall_info"
	if middleware.M.KafkaProducer != nil {
		_, _, err := middleware.M.KafkaProducer.SendMessage(
			&sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(info),
			})
		if err != nil {
			ilog.Log.Warnf("Failed to send message to Kafka: %s", err)
		}
		ilog.Log.Infof("Message sent to topic %s, data is %s\n", topic, info)
	}
}

func SendExploreNote2Redis(ctx context.Context, reqContext *model.NoteReqContext, noteId string) {
	if middleware.M.Conf.MiddleWareControl.IsRedis {
		userViewKey := model.UserViewRedisKeyName + model.ConcatSign + strconv.FormatUint(reqContext.UserId, 10)
		member := &redis.Z{}
		member.Member = noteId
		member.Score = float64(time.Now().Unix())
		if os.Getenv("REC_ENV") == "prd" {
			err := middleware.M.RedisClusterCli.ZAdd(ctx, userViewKey, member).Err()
			if err != nil {
				ilog.Log.Errorf("userViewKey = %s, has %s\n", userViewKey, err)
			}
		} else {
			err := middleware.M.RedisCli.ZAdd(ctx, userViewKey, member).Err()
			if err != nil {
				ilog.Log.Errorf("userViewKey = %s, has %s\n", userViewKey, err)
			}
		}
	} else {
		ilog.Log.Errorf("Redis can not open!")
	}
}

func SendExploreNoteString2Redis(ctx context.Context, reqContext *model.NoteReqContext, userViewStrs string) {
	if middleware.M.Conf.MiddleWareControl.IsRedis {
		userViewKey := model.UserViewRedisKeyName + model.ConcatSign + strconv.FormatUint(reqContext.UserId, 10)
		if os.Getenv("REC_ENV") == "prd" {
			err := middleware.M.RedisClusterCli.Append(ctx, userViewKey, userViewStrs+",").Err()
			if err != nil {
				ilog.Log.Errorf("userViewKey = %s, has %s\n", userViewKey, err)
			}
		} else {
			err := middleware.M.RedisCli.Append(ctx, userViewKey, userViewStrs+",").Err()
			if err != nil {
				ilog.Log.Errorf("userViewKey = %s, has %s\n", userViewKey, err)
			}
		}
	} else {
		ilog.Log.Errorf("Redis can not open!")
	}
}

func DeletedStringUserViewKeys(ctx context.Context, url string) {
	fs := &robot.FeishuRobot{Url: url}
	if middleware.M.Conf.MiddleWareControl.IsRedis {
		if os.Getenv("REC_ENV") == "prd" {
			userViewKeys, keysErr := middleware.M.RedisClusterCli.Keys(ctx, model.UserViewRedisKeyName+"*").Result()
			if keysErr != nil {
				ilog.Log.Errorf("Get redis all userViewKey Error:", keysErr)
				return
			}
			ukLenSum := 0
			for _, userViewKey := range userViewKeys {
				userViewKeyStr, err := middleware.M.RedisClusterCli.Get(ctx, userViewKey).Result()
				if err == nil {
					userViewKeyStrList := strings.Split(userViewKeyStr, ",")
					// 截取掉10%, 防止曝光存储过大
					ukLen := len(userViewKeyStrList)
					in := int(float32(ukLen) * 0.1)
					newUserViewList := strings.Join(userViewKeyStrList[in:], ",")
					err2 := middleware.M.RedisClusterCli.Set(ctx, userViewKey, newUserViewList, 0).Err()
					if err2 != nil {
						ilog.Log.Errorf("Set userViewKey = %s, has %s\n", userViewKey, err)
					}
					ukLenSum += ukLen
				}
			}
			fs.Send2robot(fmt.Sprintf(
				"Update all userViewKey count = %d, avg len = %d", len(userViewKeys), ukLenSum/len(userViewKeys)))
		} else {
			userViewKeys, keysErr := middleware.M.RedisCli.Keys(ctx, model.UserViewRedisKeyName+"*").Result()
			if keysErr != nil {
				ilog.Log.Errorf("Get redis all userViewKey Error:", keysErr)
				return
			}
			ukLenSum := 0
			for _, userViewKey := range userViewKeys {
				userViewKeyStr, err := middleware.M.RedisCli.Get(ctx, userViewKey).Result()
				if err == nil {
					userViewKeyStrList := strings.Split(userViewKeyStr, ",")
					// 截取掉10%, 防止曝光存储过大
					ukLen := len(userViewKeyStrList)
					in := int(float32(ukLen) * 0.1)
					newUserViewList := strings.Join(userViewKeyStrList[in:], ",")
					err2 := middleware.M.RedisCli.Set(ctx, userViewKey, newUserViewList, 0).Err()
					if err2 != nil {
						ilog.Log.Errorf("Set userViewKey = %s, has %s\n", userViewKey, err)
					}
					ukLenSum += ukLen
				}
			}
			fs.Send2robot(fmt.Sprintf(
				"Update all userViewKey count = %d, avg len = %d", len(userViewKeys), ukLenSum/len(userViewKeys)))
		}
	} else {
		ilog.Log.Errorf("Redis can not open!")
	}
}
