package model

import (
	"context"
	"gitlab.com/cher8/lion/common/abtest"
	"gitlab.com/cher8/lion/common/ilog"
	"gitlab.com/cher8/lion/common/middleware"
	"strings"
)

type NoteReqContext struct {
	// 笔记请求上下文
	ReqId   string `json:"req_id"`
	Version string `json:"version"`
	UserId  uint64 `json:"user_id"` //用户id
	Sex     string `json:"sex"`     //用户性别
	IP      string `json:"ip"`
	OS      string `json:"os"` //操作系统 安卓 ios
	IsDebug bool   `json:"debug"`
	TraceId string `json:"trace_id"`
	// 以下为推荐引擎内部生成的字段 不是入参数
	UniqueId    string                 //唯一id 用于dump
	Timestamp   int64                  //时间
	User        *User                  `json:"user"`
	NoteIds     []*Note                `json:"note_ids"`
	NoteSummary *NoteSummary           `json:"note_summary"`
	Conf        map[string]interface{} `json:"conf"`
	AbTestConf  *abtest.ExpGroup       `json:"ab_test_conf"`
}

type RecallScore struct {
	Name   string  `json:"name"`
	Len    int     `json:"len"`
	Score  float32 `json:"score"`
	Weight float32 `json:"weight"`
}

type Tags struct {
	ParentTag string `json:"parent_tag"`
	ChildTag  string `json:"sub_tag"`
	SexTag    string `json:"sex_tag"`
}

type Note struct {
	NoteId      uint64       `json:"note_id"`
	Author      string       `json:"author"`
	Tags        *Tags        `json:"tags"`
	PubTime     int64        `json:"pub_time"`
	ChannelName string       `json:"channel_name"`
	RecallScore *RecallScore `json:"recall_score"`
	MergeScore  float32      `json:"merge_score"`
	CoarseScore float32      `json:"coarse_score"`
}

type NoteSummary struct {
	NoteNewRank string `json:"note_new_rank"`
	NoteHotRank string `json:"note_hot_rank"`
}

func (ni *NoteSummary) FindNoteNewRankList(ctx context.Context) {
	if middleware.M.RedisCli != nil {
		noteNewRank, err := middleware.M.RedisCli.Get(
			ctx,
			NoteTopNewRedisKeyName).Result()
		if err == nil {
			ilog.Log.Infof("noteNewRank: %d", len(strings.Split(noteNewRank, ",")))
			ni.NoteNewRank = noteNewRank
		} else {
			ilog.Log.Infof("No find noteNewRank! %s", err)
			ni.NoteNewRank = ""
		}
	} else if middleware.M.RedisClusterCli != nil {
		noteNewRank, err := middleware.M.RedisClusterCli.Get(
			ctx,
			NoteTopNewRedisKeyName).Result()
		if err == nil {
			ilog.Log.Infof("noteNewRank: %d", len(strings.Split(noteNewRank, ",")))
			ni.NoteNewRank = noteNewRank
		} else {
			ilog.Log.Infof("No find noteNewRank! %s", err)
			ni.NoteNewRank = ""
		}
	} else {
		ni.NoteNewRank = ""
		ilog.Log.Infof("redis is nil, No noteNewRank!")
	}
}

func (ni *NoteSummary) FindNoteHotRankList(ctx context.Context) {
	if middleware.M.RedisCli != nil {
		noteHotRank, err := middleware.M.RedisCli.Get(
			ctx,
			NoteTopHotRedisKeyName).Result()
		if err == nil {
			ilog.Log.Infof("noteHotRank: %d", len(strings.Split(noteHotRank, ",")))
			ni.NoteHotRank = noteHotRank
		} else {
			ilog.Log.Infof("No find noteHotRank! %s", err)
			ni.NoteHotRank = ""
		}
	} else if middleware.M.RedisClusterCli != nil {
		noteHotRank, err := middleware.M.RedisClusterCli.Get(
			ctx,
			NoteTopHotRedisKeyName).Result()
		if err == nil {
			ilog.Log.Infof("noteHotRank: %d", len(strings.Split(noteHotRank, ",")))
			ni.NoteHotRank = noteHotRank
		} else {
			ilog.Log.Infof("No find noteHotRank! %s", err)
			ni.NoteHotRank = ""
		}
	} else {
		ni.NoteHotRank = ""
		ilog.Log.Infof("redis is nil, No noteHotRank!")
	}
}
