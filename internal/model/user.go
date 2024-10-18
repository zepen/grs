package model

import (
	"context"
	"gitlab.com/cher8/lion/common/ilog"
	"gitlab.com/cher8/lion/common/middleware"
	"gitlab.com/cher8/lion/common/util"
	"recommend-server/apis"
	"strings"
)

type UserSeqInfo struct {
	IntentList       string `json:"intent_list"`        // 用户兴趣标签集合
	ClickSeq         string `json:"click_seq"`          // 用户点击笔记序列
	QueryClickSeq    string `json:"query_click_seq"`    // 用户搜索点击笔记序列
	ViewSeq          string `json:"view_seq"`           // 用户浏览笔记序列
	LikeSeq          string `json:"like_seq"`           // 用户点赞笔记序列
	DislikeNoteSeq   string `json:"dislike_note_seq"`   // 用户不喜欢帖子集合
	DislikeAuthorSeq string `json:"dislike_author_seq"` // 用户不喜欢作者集合
	DislikeTagsSeq   string `json:"dislike_tags_seq"`   // 用户不喜欢子标签集合
	UserFollowSeq    string `json:"user_follow_seq"`    // 用户关注作者集合
	UserCFList       string `json:"user_cf_list"`       // 用户协同过滤列表
	UserFromChannel  string `json:"user_from_channel"`  // 用户来自渠道
}

type User struct {
	UserId      uint64       `json:"user_id"`
	UserSeqInfo *UserSeqInfo `json:"user_seq_info"`
}

func (u *User) FindUserIntentList(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		intentList, err := middleware.M.RedisCli.Get(
			ctx,
			UserIntentListRedisKeyName+ConcatSign+req.UserId).Result()
		intentList = strings.Replace(intentList, "\"", "", -1) // 双引号替换，java写入redis的value带有双引号
		if err == nil {
			u.UserSeqInfo.IntentList = intentList
			ilog.Log.Infof("user id = %s, intentList: %s", req.UserId, intentList)
		} else {
			u.UserSeqInfo.IntentList = util.RandGroupLabel(16, 3)
			ilog.Log.Infof("No find user intentList! %s, generate intentList = %s", err, u.UserSeqInfo.IntentList)
		}
	} else if middleware.M.RedisClusterCli != nil {
		intentList, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserIntentListRedisKeyName+ConcatSign+req.UserId).Result()
		intentList = strings.Replace(intentList, "\"", "", -1) // 双引号替换，java写入redis的value带有双引号
		if err == nil {
			u.UserSeqInfo.IntentList = intentList
			ilog.Log.Infof("user id = %s, intentList: %s", req.UserId, intentList)
		} else {
			u.UserSeqInfo.IntentList = util.RandGroupLabel(16, 3)
			ilog.Log.Infof("No find user intentList! %s, generate intentList = %s", err, u.UserSeqInfo.IntentList)
		}
	} else {
		u.UserSeqInfo.IntentList = util.RandGroupLabel(16, 3)
		ilog.Log.Infof("redis is nil, No find user intentList!, generate intentList = %s", u.UserSeqInfo.IntentList)
	}
}

func (u *User) FindUserClickSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		clickSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserClickRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.ClickSeq = clickSeq
			ilog.Log.Infof("user id = %s, clickSeq: %s", req.UserId, clickSeq)
		} else {
			u.UserSeqInfo.ClickSeq = ""
			ilog.Log.Infof("No find user clickSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		clickSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserClickRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.ClickSeq = clickSeq
			ilog.Log.Infof("user id = %s, clickSeq: %s", req.UserId, clickSeq)
		} else {
			u.UserSeqInfo.ClickSeq = ""
			ilog.Log.Infof("No find user clickSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.ClickSeq = ""
		ilog.Log.Infof("redis is nil, No clickSeq!")
	}
}

func (u *User) FindUserQueryClickSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		queryClickSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserQueryClickRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.QueryClickSeq = queryClickSeq
			ilog.Log.Infof("user id = %s, queryClickSeq: %s", req.UserId, queryClickSeq)
		} else {
			u.UserSeqInfo.QueryClickSeq = ""
			ilog.Log.Infof("No find user queryClickSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		queryClickSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserQueryClickRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.QueryClickSeq = queryClickSeq
			ilog.Log.Infof("user id = %s, queryClickSeq: %s", req.UserId, queryClickSeq)
		} else {
			u.UserSeqInfo.QueryClickSeq = ""
			ilog.Log.Infof("No find user queryClickSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.QueryClickSeq = ""
		ilog.Log.Infof("redis is nil, No queryClickSeq!")
	}
}

func (u *User) FindUserLikeSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		likeSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserLikeSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.QueryClickSeq = likeSeq
			ilog.Log.Infof("user id = %s, likeSeq: %s", req.UserId, likeSeq)
		} else {
			u.UserSeqInfo.QueryClickSeq = ""
			ilog.Log.Infof("No find user likeSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		likeSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserLikeSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.QueryClickSeq = likeSeq
			ilog.Log.Infof("user id = %s, likeSeq: %s", req.UserId, likeSeq)
		} else {
			u.UserSeqInfo.QueryClickSeq = ""
			ilog.Log.Infof("No find user likeSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.QueryClickSeq = ""
		ilog.Log.Infof("redis is nil, No likeSeq!")
	}
}

func (u *User) FindUserDisLikeNoteSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		disLikeNoteSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserDisLikeNoteSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.DislikeNoteSeq = disLikeNoteSeq
			ilog.Log.Infof("user id = %s, disLikeNoteSeq len = %d", req.UserId, len(strings.Split(disLikeNoteSeq, ",")))
		} else {
			u.UserSeqInfo.DislikeNoteSeq = ""
			ilog.Log.Infof("No find user disLikeNoteSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		disLikeNoteSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserDisLikeNoteSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.DislikeNoteSeq = disLikeNoteSeq
			ilog.Log.Infof("user id = %s, disLikeNoteSeq len = %d", req.UserId, len(strings.Split(disLikeNoteSeq, ",")))
		} else {
			u.UserSeqInfo.DislikeNoteSeq = ""
			ilog.Log.Infof("No find user disLikeNoteSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.DislikeNoteSeq = ""
		ilog.Log.Infof("redis is nil, No disLikeNoteSeq!")
	}
}

func (u *User) FindUserDisLikeAuthorSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		disLikeAuthorSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserDisLikeAuthorSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.DislikeAuthorSeq = disLikeAuthorSeq
			ilog.Log.Infof("user id = %s, disLikeAuthorSeq = %s", req.UserId, disLikeAuthorSeq)
		} else {
			u.UserSeqInfo.DislikeAuthorSeq = ""
			ilog.Log.Infof("No find user disLikeAuthorSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		disLikeAuthorSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserDisLikeAuthorSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.DislikeAuthorSeq = disLikeAuthorSeq
			ilog.Log.Infof("user id = %s, disLikeAuthorSeq = %s", req.UserId, disLikeAuthorSeq)
		} else {
			u.UserSeqInfo.DislikeAuthorSeq = ""
			ilog.Log.Infof("No find user disLikeAuthorSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.DislikeAuthorSeq = ""
		ilog.Log.Infof("redis is nil, No disLikeAuthorSeq!")
	}
}

func (u *User) FindUserDisLikeTagsSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		disLikeTagsSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserDisLikeTagsSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.DislikeTagsSeq = disLikeTagsSeq
			ilog.Log.Infof("user id = %s, disLikeTagsSeq = %s", req.UserId, disLikeTagsSeq)
		} else {
			u.UserSeqInfo.DislikeTagsSeq = ""
			ilog.Log.Infof("No find user disLikeTagsSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		disLikeTagsSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserDisLikeTagsSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.DislikeTagsSeq = disLikeTagsSeq
			ilog.Log.Infof("user id = %s, disLikeTagsSeq = %s", req.UserId, disLikeTagsSeq)
		} else {
			u.UserSeqInfo.DislikeTagsSeq = ""
			ilog.Log.Infof("No find user disLikeTagsSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.DislikeTagsSeq = ""
		ilog.Log.Infof("redis is nil, No disLikeTagsSeq!")
	}
}

func (u *User) FindUserFollowSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		userFollowSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserFollowSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.UserFollowSeq = userFollowSeq
			ilog.Log.Infof("user id = %s, userFollowSeq len = %d", req.UserId, len(strings.Split(userFollowSeq, ",")))
		} else {
			ilog.Log.Infof("No find user userFollowSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		userFollowSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserFollowSeqRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.UserFollowSeq = userFollowSeq
			ilog.Log.Infof("user id = %s, userFollowSeq len = %d", req.UserId, len(strings.Split(userFollowSeq, ",")))
		} else {
			ilog.Log.Infof("No find user userFollowSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.UserFollowSeq = ""
		ilog.Log.Infof("redis is nil, No userFollowSeq!")
	}
}

func (u *User) FindUserViewSeq(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		viewSeq, err := middleware.M.RedisCli.Get(
			ctx,
			UserViewRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.ViewSeq = viewSeq
			ilog.Log.Infof("user id = %s, viewSeq len = %d", req.UserId, len(strings.Split(viewSeq, ",")))
		} else {
			u.UserSeqInfo.ViewSeq = ""
			ilog.Log.Infof("No find user viewSeq! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		viewSeq, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserViewRedisKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.ViewSeq = viewSeq
			ilog.Log.Infof("user id = %s, viewSeq len = %d", req.UserId, len(strings.Split(viewSeq, ",")))
		} else {
			u.UserSeqInfo.ViewSeq = ""
			ilog.Log.Infof("No find user viewSeq! %s", err)
		}
	} else {
		u.UserSeqInfo.ViewSeq = ""
		ilog.Log.Infof("redis is nil, No viewSeq!")
	}
}

func (u *User) FindUserCFList(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		userCFList, err := middleware.M.RedisCli.Get(
			ctx,
			UserCFListKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.UserCFList = userCFList
			ilog.Log.Infof("user id = %s, userCFList len = %d", req.UserId, len(strings.Split(userCFList, ",")))
		} else {
			u.UserSeqInfo.UserCFList = ""
			ilog.Log.Infof("No find user userCFList! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		userCFList, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserCFListKeyName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.UserCFList = userCFList
			ilog.Log.Infof("user id = %s, userCFList len = %d", req.UserId, len(strings.Split(userCFList, ",")))
		} else {
			u.UserSeqInfo.UserCFList = ""
			ilog.Log.Infof("No find user userCFList! %s", err)
		}
	} else {
		u.UserSeqInfo.UserCFList = ""
		ilog.Log.Infof("redis is nil, No userCFList!")
	}
}

func (u *User) FindUserFromChannel(ctx context.Context, req *apis.UserRequest) {
	if middleware.M.RedisCli != nil {
		userFromChannel, err := middleware.M.RedisCli.Get(
			ctx,
			UserFromChannelName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.UserFromChannel = userFromChannel
			ilog.Log.Infof("user id = %s, userFromChannel = %s", req.UserId, u.UserSeqInfo.UserFromChannel)
		} else {
			u.UserSeqInfo.UserFromChannel = ""
			ilog.Log.Infof("No find user userFromChannel! %s", err)
		}
	} else if middleware.M.RedisClusterCli != nil {
		userFromChannel, err := middleware.M.RedisClusterCli.Get(
			ctx,
			UserFromChannelName+ConcatSign+req.UserId).Result()
		if err == nil {
			u.UserSeqInfo.UserFromChannel = userFromChannel
			ilog.Log.Infof("user id = %s, userFromChannel = %s", req.UserId, u.UserSeqInfo.UserFromChannel)
		} else {
			u.UserSeqInfo.UserFromChannel = ""
			ilog.Log.Infof("No find user userFromChannel! %s", err)
		}
	} else {
		u.UserSeqInfo.UserFromChannel = ""
		ilog.Log.Infof("redis is nil, No userFromChannel!")
	}
}
